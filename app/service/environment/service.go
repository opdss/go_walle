package environment

import (
	"errors"
	"go-walle/app/model"
	"go-walle/app/service/common"
	"gorm.io/gorm"
	"sync"
)

var (
	service     *Service
	onceService sync.Once
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	onceService.Do(func() {
		service = &Service{db: db}
	})
	return service
}

func (srv *Service) List(params *ListReq) (total int64, list []*model.Environment, err error) {
	_db := srv.db.Model(&model.Environment{}).Where(model.Environment{SpaceId: params.SpaceId})
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).Preload("Space").Find(&list).Error
	return
}

func (srv *Service) Create(params *CreateReq) error {
	return srv.db.Create(&model.Environment{
		SpaceId:     params.SpaceId,
		Name:        params.Name,
		Description: params.Description,
		Status:      params.Status,
		Color:       params.Color,
	}).Error
}

func (srv *Service) Update(params *UpdateReq) error {
	return srv.db.Model(model.Environment{}).
		Select(params.Fields()).
		Where(model.Environment{SpaceId: params.SpaceId, ID: params.ID}).
		Updates(params).Error
}

// Delete 环境下必须没有项目了才能删除
func (srv *Service) Delete(spaceWithId *common.SpaceWithId) error {
	total := srv.db.Model(&model.Environment{ID: spaceWithId.ID}).Association("Projects").Count()
	if total > 0 {
		return errors.New("该环境还存在项目，不允许删除，如需要删除，请先删除该环境下所有项目")
	}
	return srv.db.Delete(&model.Environment{SpaceId: spaceWithId.SpaceId, ID: spaceWithId.ID}).Error
}

func (srv *Service) Detail(spaceWithId *common.SpaceWithId) (m *model.Environment, err error) {
	err = srv.db.Where(spaceWithId).First(&m).Error
	return
}
