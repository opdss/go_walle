package member

import (
	"github.com/wuzfei/go-helper/slices"
	"github.com/zeebo/errs"
	"go-walle/app/global"
	"go-walle/app/model"
	"go-walle/app/service/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	service     *Service
	onceService sync.Once
)

type Service struct {
	db *gorm.DB
}

func NewService() *Service {
	onceService.Do(func() {
		service = &Service{db: global.DB}
	})
	return service
}

func (srv *Service) Store(params StoreReq) (err error) {
	user := model.User{}
	err = srv.db.First(&user, params.UserId).Error
	if err != nil {
		return
	}
	if !user.Status.IsEnable() {
		return errs.New("该用户已被禁用")
	}
	member := model.Member{
		SpaceId: params.SpaceId,
		UserId:  params.UserId,
		Role:    params.Role,
	}
	return srv.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "space_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Create(&member).Error
}

func (srv *Service) Delete(spaceAndId *common.SpaceWithId) (err error) {
	return srv.db.Where(spaceAndId).Delete(&model.Member{}).Error
}

func (srv *Service) List(params ListReq) (total int64, res []*ListItem, err error) {
	members := make([]*model.Member, 0)
	_db := srv.db.Model(&model.Member{}).Where("space_id= ?", params.SpaceId)
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).Preload("User").Find(&members).Error
	if err != nil {
		return
	}
	res = slices.Map(members, func(item *model.Member, k int) *ListItem {
		return &ListItem{
			SpaceId:   item.SpaceId,
			UserId:    item.UserId,
			Username:  item.User.Username,
			Email:     item.User.Email,
			Role:      item.Role,
			Status:    item.User.Status,
			CreatedAt: item.CreatedAt,
		}
	})
	return
}
