package deploy

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/wuzfei/go-helper/slices"
	"go-walle/app/global"
	"go-walle/app/internal/errcode"
	"go-walle/app/model"
	"go-walle/app/service/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	service     *Service
	onceService sync.Once
)

type Service struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewService() *Service {
	onceService.Do(func() {
		service = &Service{
			db:  global.DB,
			log: global.Log,
		}
	})
	return service
}

func (srv *Service) List(params *ListReq) (total int64, list []*model.Task, err error) {
	_db := srv.db.Model(&model.Task{}).Where("space_id=?", params.SpaceId)
	err = _db.Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).
		Preload("User").
		Preload("Project").
		Preload("Environment").
		Order("id desc").
		Find(&list).Error
	return
}

// Create 创建上线单
func (srv *Service) Create(params *CreateReq) error {
	project := &model.Project{SpaceId: params.SpaceId, ID: params.ProjectId}
	err := srv.db.Model(&project).Where(project).Preload("Environment").Preload("Servers").First(&project).Error
	if err != nil {
		return err
	}
	if !project.Status.IsEnable() || !project.Environment.Status.IsEnable() {
		return errors.New("该项目或者该环境暂停上线，请联系相关负责人")
	}
	serverIds := slices.Map(project.Servers, func(item model.Server, k int) int64 {
		return item.ID
	})
	m := &model.Task{
		Name:          params.Name,
		SpaceId:       params.SpaceId,
		UserId:        params.UserId,
		ProjectId:     project.ID,
		EnvironmentId: project.Environment.ID,
		Tag:           params.Tag,
		Branch:        params.Branch,
		CommitId:      params.CommitId,
		ServerIds:     slices.Intersect(serverIds, params.ServerIds),
	}
	m.Status = model.TaskStatusAudit
	if project.TaskAudit == 1 {
		m.Status = model.TaskStatusWaiting
	}
	if len(m.ServerIds) == 0 {
		return errcode.ErrRequest.Wrap(errors.New("服务器选择错误"))
	}
	return srv.db.Create(m).Error
}

// Detail 上线单详情
func (srv *Service) Detail(spaceAndId *common.SpaceWithId) (taskDetail *model.Task, err error) {
	taskDetail = &model.Task{}
	err = srv.db.Where(spaceAndId).
		Preload("Project").
		First(&taskDetail).
		Error
	if err != nil {
		return
	}
	if len(taskDetail.ServerIds) > 0 {
		servers := make([]*model.Server, 0)
		err = srv.db.Where("id in ?", []int64(taskDetail.ServerIds)).Find(&servers).Error
		if err != nil {
			return
		}
		taskDetail.Servers = servers
	}
	return
}

// Delete 删除
func (srv *Service) Delete(spaceId int64) (m *model.Space, err error) {
	err = srv.db.First(&m, spaceId).Error
	return
}

// Audit 审核
func (srv *Service) Audit(params *AuditReq) (err error) {
	var m *model.Task
	err = srv.db.Where("space_id = ? and id = ?", params.SpaceId, params.ID).First(&m).Error
	if err != nil {
		return
	}
	if m.Status != model.TaskStatusWaiting {
		return errors.New("审核失败，该上线单并未处理待审核状态")
	}

	m.AuditUserId = params.AuditUserId
	if params.Audit {
		m.Status = model.TaskStatusAudit
	} else {
		m.Status = model.TaskStatusReject
	}
	return srv.db.Select("status", "audit_user_id").Updates(&m).Error
}

// Release 发布
func (srv *Service) Release(spaceAndId *common.SpaceWithId, userId int64) (err error) {
	//上线单详情
	taskDetail, err := srv.getTask(spaceAndId, "Project", "Environment")
	if err != nil {
		return
	}
	if len(taskDetail.ServerIds) > 0 {
		servers := make([]*model.Server, 0)
		err = srv.db.Find(&servers, []int64(taskDetail.ServerIds)).Error
		if err != nil {
			return
		}
		taskDetail.Servers = servers
	}
	deployTask, err := CreateDeployTask(taskDetail, userId)
	if err != nil {
		return err
	}
	return deployTask.Start()
}

// Stop 停止发布
func (srv *Service) Stop(spaceId, id, userId int64) (err error) {
	//上线单详情
	deployTask := GetDeployTask(id)
	if deployTask == nil {
		return errors.New("已经终止发布")
	}
	return deployTask.Stop()
}

// Rollback 回滚
func (srv *Service) Rollback(spaceId int64) (m *model.Space, err error) {
	err = srv.db.First(&m, spaceId).Error
	return
}

// Console 回滚
func (srv *Service) ConsoleNew(wsConn *websocket.Conn, spaceAndId *common.SpaceWithId) {
	var err error
	defer func() {
		if err != nil {
			srv.log.Error("获取发布日志出错",
				zap.Int64("spaceId", spaceAndId.SpaceId),
				zap.Int64("taskId", spaceAndId.ID),
				zap.Error(err))
		}
	}()
	var taskModel *model.Task
	taskModel, err = srv.getTask(spaceAndId)
	if err != nil {
		return
	}
	task := GetDeployTask(taskModel.ID)
	recordChan := make(chan *TaskConsoleMsg, 1)
	defer close(recordChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	writeErr := make(chan error)
	defer close(writeErr)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-recordChan:
				str, _ := json.Marshal(msg)
				_err := wsConn.WriteMessage(websocket.TextMessage, str)
				if _err != nil {
					global.Log.Error("ws发送失败", zap.Error(err), zap.Int64("id", msg.Records[0].ID))
					writeErr <- _err
					return
				}
				writeErr <- nil
			}
		}
	}()

	lastId := int64(0)
	res := make([]*model.Record, 0)
	if task != nil {
		getRecords := func() {
			if err = srv.db.Where("task_id = ?", taskModel.ID).Where("id > ?", lastId).Order("created_at asc").Find(&res).Error; err != nil {
				return
			}
			for _, v := range res {
				recordChan <- &TaskConsoleMsg{
					Type:    TaskConsoleMsgAppend,
					Records: []*model.Record{v},
				}
				lastId = v.ID
				if err = <-writeErr; err != nil {
					return
				}
			}
		}
		for {
			select {
			case <-task.IsStop():
				getRecords()
				global.Log.Debug("发布任务已经完成", zap.Int64("taskId", spaceAndId.ID))
				return
			default:
				getRecords()
				time.Sleep(time.Second * 1)
			}
		}
	} else {
		if err = srv.db.Where("task_id = ?", taskModel.ID).Order("created_at asc").Find(&res).Error; err != nil {
			return
		}
		for _, v := range res {
			recordChan <- &TaskConsoleMsg{
				Type:    TaskConsoleMsgAppend,
				Records: []*model.Record{v},
			}
			if err = <-writeErr; err != nil {
				return
			}
		}
	}

}

// Console 回滚
func (srv *Service) Console(wsConn *websocket.Conn, spaceAndId *common.SpaceWithId) {
	var err error
	defer func() {
		if err != nil {
			srv.log.Error("获取发布日志出错",
				zap.Int64("spaceId", spaceAndId.SpaceId),
				zap.Int64("taskId", spaceAndId.ID),
				zap.Error(err))
		}
	}()
	var taskModel *model.Task
	taskModel, err = srv.getTask(spaceAndId)
	if err != nil {
		return
	}
	task := GetDeployTask(taskModel.ID)
	recordChan := make(chan *TaskConsoleMsg, 1)
	defer close(recordChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	writeErr := make(chan error)
	defer close(writeErr)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-recordChan:
				str, _ := json.Marshal(msg)
				_err := wsConn.WriteMessage(websocket.TextMessage, str)
				if _err != nil {
					global.Log.Error("ws发送失败", zap.Error(err), zap.Int64("id", msg.Records[0].ID))
					writeErr <- _err
					return
				}
				writeErr <- nil
			}
		}
	}()

	lastId := int64(0)
	res := make([]*model.Record, 0)
	if task != nil {
		getRecords := func() {
			if err = srv.db.Where("task_id = ?", taskModel.ID).Where("id > ?", lastId).Order("created_at asc").Find(&res).Error; err != nil {
				return
			}
			for _, v := range res {
				recordChan <- &TaskConsoleMsg{
					Type:    TaskConsoleMsgAppend,
					Records: []*model.Record{v},
				}
				lastId = v.ID
				if err = <-writeErr; err != nil {
					return
				}
			}
		}
		for {
			select {
			case <-task.IsStop():
				getRecords()
				global.Log.Debug("发布任务已经完成", zap.Int64("taskId", spaceAndId.ID))
				return
			default:
				getRecords()
				time.Sleep(time.Second * 1)
			}
		}
	} else {
		if err = srv.db.Where("task_id = ?", taskModel.ID).Order("created_at asc").Find(&res).Error; err != nil {
			return
		}
		for _, v := range res {
			recordChan <- &TaskConsoleMsg{
				Type:    TaskConsoleMsgAppend,
				Records: []*model.Record{v},
			}
			if err = <-writeErr; err != nil {
				return
			}
		}
	}

}

func (srv *Service) getTask(spaceAndId *common.SpaceWithId, preloads ...string) (*model.Task, error) {
	//上线单详情
	taskDetail := model.Task{}
	_db := srv.db.Where(spaceAndId)
	for _, pre := range preloads {
		_db = _db.Preload(pre)
	}
	err := _db.First(&taskDetail).Error
	return &taskDetail, err
}
