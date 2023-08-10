package user

import (
	"errors"
	"github.com/wuzfei/go-helper/slices"
	"go-walle/app/internal/constants"
	"go-walle/app/internal/errcode"
	"go-walle/app/model"
	"go-walle/app/pkg/jwt"
	"go-walle/app/service/common"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	jwt *jwt.Jwt
}

func NewService(log *zap.Logger, db *gorm.DB, jwt *jwt.Jwt) *Service {
	onceService.Do(func() {
		service = &Service{
			log: log,
			db:  db,
			jwt: jwt,
		}
	})
	return service
}

// Login 登陆
func (srv *Service) Login(params *LoginReq) (*LoginRes, error) {
	m := model.User{}
	err := srv.db.Where("email = ?", params.Email).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errcode.ErrInvalidPwd
		}
		return nil, err
	}
	if m.Status.IsDisable() {
		return nil, errcode.ErrUserDisabled
	}
	if !(bcrypt.CompareHashAndPassword(m.Password, []byte(params.Password)) == nil) {
		return nil, errcode.ErrInvalidPwd
	}
	//生成token
	res := LoginRes{}
	res.Token, res.TokenExpire, err = srv.jwt.CreateToken(jwt.TokenPayload{
		UserId:   m.ID,
		Email:    m.Email,
		Username: m.Username,
	})
	if err != nil {
		return nil, err
	}
	res.UserId = m.ID
	//记住登陆
	if params.Remember {
		res.RefreshToken, res.RefreshTokenExpire, err = srv.jwt.CreateRefreshToken(jwt.TokenPayload{
			UserId:    m.ID,
			Username:  m.Username,
			IsRefresh: true,
		})
		if err != nil {
			return nil, err
		}
		m.RememberToken = res.RefreshToken
		if err = srv.db.Select("remember_token").Updates(&m).Error; err != nil {
			return nil, err
		}
	}
	return &res, nil
}

// RefreshToken 刷新token
func (srv *Service) RefreshToken(params *RefreshTokenReq) (res *LoginRes, err error) {
	jwtClaims, err := srv.jwt.ValidateToken(params.RefreshToken)
	if err != nil {
		return
	}
	m := model.User{}
	err = srv.db.First(&m, jwtClaims.UserId).Error
	if err != nil {
		return
	}
	if m.Status != constants.StatusEnable {
		return nil, errcode.ErrUserDisabled
	}
	if m.RememberToken != params.RefreshToken {
		return nil, errcode.ErrInvalidParams.New("refresh token 错误")
	}

	res = &LoginRes{}
	res.Token, res.TokenExpire, err = srv.jwt.CreateToken(jwt.TokenPayload{
		UserId:   m.ID,
		Username: m.Username,
	})
	if err != nil {
		return
	}
	res.UserId = m.ID
	res.RefreshToken, res.RefreshTokenExpire, err = srv.jwt.CreateRefreshToken(jwt.TokenPayload{
		UserId:    m.ID,
		Username:  m.Username,
		IsRefresh: true,
	})
	if err != nil {
		return nil, err
	}
	m.RememberToken = res.RefreshToken
	if err = srv.db.Select("remember_token").Updates(&m).Error; err != nil {
		return nil, err
	}
	return
}

// Logout 退出
func (srv *Service) Logout(userId int64) (err error) {
	m := model.User{}
	err = srv.db.First(&m, userId).Error
	if err != nil {
		return
	}
	m.RememberToken = ""
	return srv.db.Select("remember_token").Updates(&m).Error
}

// Create 创建新用户
func (srv *Service) Create(params *CreateReq) (err error) {
	m := model.User{}
	var exists int64
	err = srv.db.Model(&m).Where("email = ?", params.Email).Count(&exists).Error
	if err != nil {
		return
	}
	if exists != 0 {
		return errcode.ErrInvalidParams.Wrap(errors.New("该用户email已存在"))
	}
	m.Username = params.Username
	m.Email = params.Email
	m.Status = params.Status

	_pwd, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	m.Password = _pwd
	return srv.db.Create(&m).Error
}

// Update 更新用户
func (srv *Service) Update(params *UpdateReq) (err error) {
	m := model.User{}
	err = srv.db.First(&m, params.ID).Error
	if err != nil {
		return
	}
	if m.ID == 0 {
		return errors.New("用户不存在")
	}
	m.ID = params.ID
	m.Username = params.Username
	m.Email = params.Email
	m.Status = params.Status
	if params.Password != "" {
		v, _err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if _err != nil {
			return _err
		}
		m.Password = v
	}
	return srv.db.UpdateColumns(&m).Error
}

// Delete 删除用户
func (srv *Service) Delete(id int64) (err error) {
	if constants.IsSuperUser(id) {
		return errors.New("超级管理员不允许删除")
	}
	return srv.db.Delete(&model.User{}, id).Error
}

// List 获取列表
func (srv *Service) List(params *ListReq) (total int64, res []*model.User, err error) {
	db := srv.db.Model(&model.User{})
	if params.Keyword != "" {
		_k := "%" + params.Keyword + "%"
		db.Where("username like ? or email like ?", _k, _k)
	}
	err = db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = db.Scopes(params.PageQuery()).Find(&res).Error
	return
}

func (srv *Service) Members(spaceId int64, params *ListReq) (total int64, res []any, err error) {
	var result []struct {
	}
	_db := srv.db.Model(&model.Member{}).Where(model.Member{SpaceId: spaceId})
	err = _db.Count(&total).Error
	if err != nil || total == 0 {
		return
	}
	err = _db.Scopes(params.PageQuery()).
		Joins("User").
		Scan(&result).
		Error
	return
}

// UserInfo 获取用户信息
func (srv *Service) UserInfo(params *common.SpaceWithId) (userInfo *GetUserInfoRes, err error) {
	m := model.User{}
	err = srv.db.First(&m, params.ID).Error
	if err != nil {
		return
	}
	role := ""
	if constants.IsSuperUser(params.ID) {
		role = string(constants.RoleSuper)
	}
	currentSpaceId := params.SpaceId
	//获取所属空间
	spaceItems, err := srv.SpacesItems(m.ID)
	if err != nil {
		return
	}
	//根据当前空间，获取当前空间id和角色
	currSpaceItem := spaceItems.Default(currentSpaceId)
	if currSpaceItem != nil {
		currentSpaceId = currSpaceItem.SpaceId
		if role == "" {
			role = currSpaceItem.Role
		}
	} else {
		currentSpaceId = 0
	}

	userInfo = &GetUserInfoRes{
		UserID:         m.ID,
		Email:          m.Email,
		Username:       m.Username,
		Role:           role,
		Status:         m.Status,
		CurrentSpaceId: currentSpaceId,
		Spaces:         spaceItems,
	}
	return
}

func (srv *Service) SpacesItems(userId int64) (SpaceItems, error) {
	spaceItems := make(SpaceItems, 0)
	if !constants.IsSuperUser(userId) {
		res, err := srv.Spaces(userId)
		if err != nil {
			return spaceItems, err
		}
		spaceItems = slices.Map(res, func(item *model.Member, k int) *SpaceItem {
			return &SpaceItem{
				SpaceId:   item.Space.ID,
				SpaceName: item.Space.Name,
				Status:    item.Space.Status,
				Role:      item.Role,
			}
		})
	} else {
		//超级管理员的处理
		var res []*model.Space
		err := srv.db.Find(&res).Error
		if err != nil {
			return spaceItems, err
		}
		spaceItems = slices.Map(res, func(item *model.Space, k int) *SpaceItem {
			return &SpaceItem{
				SpaceId:   item.ID,
				SpaceName: item.Name,
				Status:    item.Status,
				Role:      string(constants.RoleSuper),
			}
		})
	}
	return spaceItems, nil
}

// Spaces 获取一个用户所有空间信息
func (srv *Service) Spaces(userId int64) (res []*model.Member, err error) {
	err = srv.db.Where("user_id = ? and status = ?", userId, 1).Preload("Space").Find(&res).Error
	if err != nil {
		return
	}
	res = slices.FilterFunc(res, func(item *model.Member) bool {
		return item.Space.ID != 0
	})
	return
}

// SpaceById 获取一个用户所有空间信息
func (srv *Service) SpaceById(userId int64, spaceId int64) (res *model.Member, err error) {
	err = srv.db.Where(model.Member{SpaceId: spaceId, UserId: userId}).First(&res).Error
	return
}

func (srv *Service) Detail(userId int64) (m *model.User, err error) {
	err = srv.db.First(&m, userId).Error
	return
}
