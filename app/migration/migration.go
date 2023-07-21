package migration

import (
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

type Migration struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewMigration(log *zap.Logger, db *gorm.DB) *Migration {
	return &Migration{
		db:  db,
		log: log,
	}
}

func (m *Migration) Setup() error {
	err := m.createTables()
	if err != nil {
		return err
	}
	return m.initUser()
}

func (m *Migration) createTables() error {
	err := m.db.AutoMigrate(
		&model.User{},
		&model.Server{},
		&model.Space{},
		&model.Project{},
		&model.Environment{},
		&model.Member{},
		&model.Record{},
		&model.Task{},
	)
	if err != nil {
		m.log.Error("创建表失败", zap.Error(err))
	}
	return err
}

func (m *Migration) initUser() error {
	defPwd := []byte(rand.StringN(12))
	_pwd, _ := bcrypt.GenerateFromPassword(defPwd, bcrypt.DefaultCost)
	mUser := model.User{
		ID:       constants.SuperUserId,
		Username: "admin",
		Email:    "admin@qq.com",
		Password: string(_pwd),
		Status:   field.StatusEnable,
	}
	err := m.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&mUser).Error
	if err != nil {
		m.log.Error("初始化admin账户失败", zap.Error(err))
		return err
	}
	m.log.Info(fmt.Sprintf("初始化admin账户成功，账号：%s, 密码：%s,请及时修改", mUser.Email, defPwd))
	return nil
}
