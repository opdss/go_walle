package space

import (
	"errors"
	"go-walle/app/model"
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
		service = &Service{
			db: db,
		}
	})
	return service
}

func (srv *Service) List(params *ListReq) (total int64, list []*model.Space, err error) {
	err = srv.db.Model(&model.Space{}).Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = srv.db.Model(&model.Space{}).Scopes(params.PageQuery()).Preload("User").Find(&list).Error
	return
}

func (srv *Service) Create(params *CreateReq) error {
	return srv.db.Create(&model.Space{
		UserId: params.UserId,
		Name:   params.Name,
		Status: params.Status,
	}).Error
}

func (srv *Service) Update(params *UpdateReq) error {
	return srv.db.Select(params.Fields()).Updates(model.Space{
		ID:     params.ID,
		Name:   params.Name,
		Status: params.Status,
		UserId: params.UserId,
	}).Error
}

// Delete 空间必须没有绑定项目才能删除
func (srv *Service) Delete(id int64) error {
	total := srv.db.Model(&model.Space{ID: id}).Association("Projects").Count()
	if total > 0 {
		return errors.New("该空间存在项目，不允许删除，如需要删除，先删除该空间下所有项目")
	}
	return srv.db.Delete(&model.Space{ID: id}).Error
}

func (srv *Service) Detail(id int64) (m *model.Space, err error) {
	err = srv.db.First(&m, id).Error
	return
}
