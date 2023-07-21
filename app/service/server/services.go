package server

import (
	"errors"
	"fmt"
	"go-walle/app/model"
	"go-walle/app/model/field"
	"go-walle/app/pkg/ssh"
	"go-walle/app/service/common"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var (
	service     *Service
	onceService sync.Once
)

type Service struct {
	log *zap.Logger
	db  *gorm.DB
	ssh *ssh.Ssh
}

func NewService(log *zap.Logger, db *gorm.DB, ssh *ssh.Ssh) *Service {
	onceService.Do(func() {
		service = &Service{db: db, log: log, ssh: ssh}
	})
	return service
}

func (srv *Service) List(params *ListReq) (total int64, list []*model.Server, err error) {
	_db := srv.db.Model(&model.Server{}).Where("space_id = ? ", params.SpaceId)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).Find(&list).Error
	return
}

func (srv *Service) Create(params *CreateReq) error {
	m := &model.Server{
		SpaceId:     params.SpaceId,
		Name:        params.Name,
		Host:        params.Host,
		Port:        params.Port,
		User:        params.User,
		Status:      field.StatusDisable,
		Description: params.Description,
	}
	_m, err := srv.FindByHostIp(m.SpaceId, m.User, m.Host, m.Port)
	if err != nil {
		return err
	}
	if _m.ID != 0 {
		return errors.New(fmt.Sprintf("已存在该主机：[%s@%s:%d]", m.User, m.Host, m.Port))
	}
	return srv.db.Create(m).Error
}

func (srv *Service) Update(params *UpdateReq) error {
	_m, err := srv.FindByHostIp(params.SpaceId, params.User, params.Host, params.Port)
	if err != nil {
		return err
	}
	if _m.ID != 0 && _m.ID != params.ID {
		return errors.New("更新错误")
	}
	return srv.db.Model(model.Server{}).Select(params.Fields()).Where(model.Server{SpaceId: params.SpaceId, ID: params.ID}).Updates(params).Error
}

func (srv *Service) Delete(spaceWith *common.SpaceWithId) error {
	return srv.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.Server{ID: spaceWith.ID}).Association("Projects").Clear()
		if err != nil {
			return err
		}
		result := tx.Where(spaceWith).Delete(&model.Server{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("删除失败")
		}
		return nil
	})
}

// FindByHostIp aa
func (srv *Service) FindByHostIp(spaceId int64, user, host string, port int) (m model.Server, err error) {
	err = srv.db.Where(&model.Server{SpaceId: spaceId, User: user, Host: host, Port: port}).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
