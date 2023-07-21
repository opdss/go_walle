package project

import (
	"errors"
	"fmt"
	"go-walle/app/model"
	"go-walle/app/model/field"
	"go-walle/app/pkg/repo"
	"go-walle/app/pkg/ssh"
	"go-walle/app/service/common"
	"gorm.io/gorm"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	service     *Service
	onceService sync.Once
)

type Service struct {
	db   *gorm.DB
	ssh  *ssh.Ssh
	repo *repo.Repos
}

func NewService(db *gorm.DB, ssh *ssh.Ssh, repo *repo.Repos) *Service {
	onceService.Do(func() {
		service = &Service{
			db:   db,
			ssh:  ssh,
			repo: repo,
		}
	})
	return service
}

func (srv *Service) List(params *ListReq) (total int64, list []*model.Project, err error) {
	where := model.Project{SpaceId: params.SpaceId}
	if params.EnvironmentId > 0 {
		where.EnvironmentId = params.EnvironmentId
	}
	_db := srv.db.Model(&where).Where(where)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).
		Preload("Space").
		Preload("Environment").Find(&list).Error
	return
}

func (srv *Service) Create(params *CreateReq) error {
	m := &model.Project{
		SpaceId: params.SpaceId,

		Name:          params.Name,
		EnvironmentId: params.EnvironmentId,
		RepoUrl:       params.RepoUrl,
		RepoMode:      params.RepoMode,
		RepoType:      params.RepoType,
		TaskAudit:     params.TaskAudit,
		Description:   params.Description,

		TargetRoot:     params.TargetRoot,
		TargetReleases: params.TargetReleases,
		KeepVersionNum: params.KeepVersionNum,

		Excludes:    params.Excludes,
		IsInclude:   params.IsInclude,
		TaskVars:    params.TaskVars,
		PrevDeploy:  params.PrevRelease,
		PostDeploy:  params.PostDeploy,
		PrevRelease: params.PrevRelease,
		PostRelease: params.PostRelease,
		Status:      field.StatusEnable,
	}
	servers := make([]model.Server, 0)
	return srv.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("space_id = ? and id in ?", params.SpaceId, params.ServerIds).Find(&servers).Error
		if err != nil {
			return err
		}
		m.Servers = servers
		return tx.Create(m).Error
	})
}

func (srv *Service) Update(params *UpdateReq) error {
	m := model.Project{}
	err := srv.db.Where("space_id = ? and id = ?", params.SpaceId, params.ID).First(&m).Error
	if err != nil {
		return err
	}
	m = model.Project{
		ID:      params.ID,
		SpaceId: params.SpaceId,

		Name:          params.Name,
		EnvironmentId: params.EnvironmentId,
		RepoUrl:       params.RepoUrl,
		RepoMode:      params.RepoMode,
		RepoType:      params.RepoType,
		TaskAudit:     params.TaskAudit,
		Description:   params.Description,

		TargetRoot:     params.TargetRoot,
		TargetReleases: params.TargetReleases,
		KeepVersionNum: params.KeepVersionNum,

		Excludes:    params.Excludes,
		IsInclude:   params.IsInclude,
		TaskVars:    params.TaskVars,
		PrevDeploy:  params.PrevRelease,
		PostDeploy:  params.PostDeploy,
		PrevRelease: params.PrevRelease,
		PostRelease: params.PostRelease,
	}
	return srv.db.Transaction(func(tx *gorm.DB) error {
		servers := make([]model.Server, 0)
		err = tx.Where("space_id = ? and id in ?", params.SpaceId, params.ServerIds).Find(&servers).Error
		if err != nil {
			return err
		}
		//晴空关联
		err = tx.Model(&model.Project{ID: params.ID}).Association("Servers").Clear()
		if err != nil {
			return err
		}
		//更新关联
		m.Servers = servers
		return tx.Model(&m).Where("space_id = ? and id = ?", params.SpaceId, params.ID).Select(params.Fields(), "Servers").UpdateColumns(m).Error
	})
}

func (srv *Service) Delete(spaceAndId *common.SpaceWithId) error {
	return srv.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Project{ID: spaceAndId.ID}).Association("Servers").Clear(); err != nil {
			return err
		}
		result := tx.Where(spaceAndId).Delete(&model.Project{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("删除失败")
		}
		return nil
	})
}

func (srv *Service) Detail(spaceAndId *common.SpaceWithId) (res model.Project, err error) {
	err = srv.db.Where(spaceAndId).Preload("Servers").First(&res).Error
	return
}

// Detection 项目检测
func (srv *Service) Detection(spaceWithId *common.SpaceWithId) (ret []*DetectionMsg, err error) {
	project := model.Project{}
	err = srv.db.Where(spaceWithId).Preload("Servers").First(&project).Error
	if err != nil {
		return
	}
	_, err = srv.repo.New(repo.TypeRepo(project.RepoType), project.RepoUrl, strconv.Itoa(int(project.ID)))
	if err != nil {
		ret = append(ret, &DetectionMsg{
			Title: "代码clone失败",
			Error: err.Error(),
			Todo:  "1、请检查仓库地址：" + project.RepoUrl + "是否正确；\n 2、请检查" + project.RepoType + "相关配置是否正确",
		})
		err = nil
	}
	if len(project.Servers) == 0 {
		ret = append(ret, &DetectionMsg{
			Title: "项目未绑定发布服务器",
			Error: "",
			Todo:  "请添加发布服务器后，在修改项目重新选择绑定",
		})
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(project.Servers))
	detectionMsgChan := make(chan *DetectionMsg, len(project.Servers))
	overChan := make(chan struct{})
	defer func() {
		close(detectionMsgChan)
		close(overChan)
	}()
	go func() {
		for msg := range detectionMsgChan {
			ret = append(ret, msg)
		}
		overChan <- struct{}{}
	}()
	for _, server := range project.Servers {
		go func(server model.Server) {
			defer wg.Done()
			re, _err := srv.ssh.NewRemoteExec(ssh.ServerConfig{User: server.User, Port: server.Port, Host: server.Host})
			if _err != nil {
				detectionMsgChan <- &DetectionMsg{
					Title: "远程目标机器免密码登录失败",
					Error: _err.Error(),
					Todo:  fmt.Sprintf("在宿主机中配置免密码登录，把宿主机用户[%s]的~/.ssh/id_rsa.pub添加到远程目标机器用户[%s]的~/.ssh/authorized_keys", server.User, server.User),
				}
				return
			}
			defer re.Close()

			webroot := filepath.Dir(project.TargetRoot)
			cmd := fmt.Sprintf("[ -d %s ] || mkdir -p %s", webroot, webroot)
			_, _err = re.Run(cmd)
			if _err != nil {
				detectionMsgChan <- &DetectionMsg{
					Title: "[" + server.Hostname() + "]远程目标机器创建目录失败",
					Error: _err.Error(),
					Todo:  fmt.Sprintf("请检查远程目标服务器用户[%s]的权限，失败执行命令：%s", server.User, cmd),
				}
				return
			}

			cmd = fmt.Sprintf("[ -L \"%s\" ] && echo \"true\" || echo \"false\"", project.TargetRoot)
			output, _err := re.Run(cmd)
			if _err != nil {
				detectionMsgChan <- &DetectionMsg{
					Title: "[" + server.Hostname() + "]目标机器执行命令失败",
					Error: _err.Error(),
					Todo:  fmt.Sprintf("请检查远程目标服务器用户[%s]的权限，失败执行命令：%s", server.User, cmd),
				}
				return
			}
			if string(output) == "false" {
				detectionMsgChan <- &DetectionMsg{
					Title: "[" + server.Hostname() + "]远程目标机器webroot不能是已建好的目录",
					Error: "远程目标机器%s webroot不能是已存在的目录，必须为软链接，你不必新建，walle会自行创建。",
					Todo:  "手工删除远程目标机器：" + server.Host + " webroot目录：" + project.TargetRoot,
				}
				return
			}
			return
		}(server)
	}
	wg.Wait()
	<-overChan
	return
}

func (srv *Service) GetBranches(spaceWithId *common.SpaceWithId) (res []repo.Branch, err error) {
	var rep repo.Repo
	rep, err = srv.getRepoBySpaceWithId(spaceWithId)
	if err != nil {
		return
	}
	return rep.Branches()
}

func (srv *Service) GetTags(spaceWithId *common.SpaceWithId) (res []repo.Tag, err error) {
	var rep repo.Repo
	rep, err = srv.getRepoBySpaceWithId(spaceWithId)
	if err != nil {
		return
	}
	return rep.Tags()
}

func (srv *Service) GetCommits(spaceWithId *common.SpaceWithId, branch string) (res []repo.Commit, err error) {
	var rep repo.Repo
	rep, err = srv.getRepoBySpaceWithId(spaceWithId)
	if err != nil {
		return
	}
	return rep.Commits(branch)
}

func (srv *Service) getRepoBySpaceWithId(spaceWithId *common.SpaceWithId) (rep repo.Repo, err error) {
	var projectModel *model.Project
	err = srv.db.Where(spaceWithId).First(&projectModel).Error
	if err != nil {
		return nil, err
	}
	if !projectModel.Status.IsEnable() {
		return nil, errors.New("该项目已经禁用")
	}
	return srv.repo.New(repo.TypeRepo(projectModel.RepoType), projectModel.RepoUrl, strconv.Itoa(int(projectModel.ID)))
}
