package migration

import (
	"errors"
	"fmt"
	"github.com/wuzfei/go-helper/rand"
	"go-walle/app/internal/constants"
	"go-walle/app/model"
	"go-walle/app/model/field"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Config struct {
	Account  string `help:"超级管理员名称" default:"admin"`
	Email    string `help:"超级管理员邮箱" default:"admin@gowalle.com"`
	Password string `help:"超级管理员密码" default:""`
}

type Migration struct {
	config *Config
	db     *gorm.DB
	log    *zap.Logger
}

func NewMigration(conf *Config, log *zap.Logger, db *gorm.DB) *Migration {
	return &Migration{
		config: conf,
		db:     db,
		log:    log,
	}
}

func (m *Migration) Setup() error {
	err := m.createTables()
	if err != nil {
		m.log.Error("创建表失败", zap.Error(err))
		return err
	}
	return m.initAdminAccount()
}

func (m *Migration) createTables() error {
	return m.db.AutoMigrate(
		&model.User{},
		&model.Server{},
		&model.Space{},
		&model.Project{},
		&model.Environment{},
		&model.Member{},
		&model.Record{},
		&model.Task{},
	)
}

// initAdmin 初始化超管账户
func (m *Migration) initAdminAccount() error {
	mUser := model.User{
		ID:       constants.SuperUserId,
		Username: "admin",
		Email:    "admin@gowalle.com",
		Status:   field.StatusEnable,
	}
	//设置密码
	var pwd []byte
	if m.config.Password == "" {
		pwd = []byte(rand.StringN(12))
	} else {
		if len(m.config.Password) < 6 || len(m.config.Password) > 32 {
			return errors.New("密码至少设置6个字符,最长32个字符")
		}
		pwd = []byte(m.config.Password)
	}
	mUser.Password, _ = bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)

	//设置账号
	if m.config.Account != "" {
		mUser.Username = m.config.Account
	}
	//设置邮箱
	if m.config.Email != "" {
		mUser.Email = m.config.Email
	}

	err := m.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&mUser).Error
	if err != nil {
		m.log.Error("初始化admin账户失败", zap.Error(err))
		return err
	}
	m.log.Info(fmt.Sprintf("初始化[%s]账户成功，账号：%s, 密码：%s,请及时修改", mUser.Username, mUser.Email, pwd))
	return nil
}
