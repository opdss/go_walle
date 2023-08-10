package model

import (
	"go-walle/app/model/field"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`

	Username      string       `gorm:"column:username;size:100;notNull;default:'';comment:用户名" json:"username"`
	Email         string       `gorm:"column:email;uniqueIndex;size:100;notNull;comment:邮箱" json:"email"`
	Password      []byte       `gorm:"column:password;size:200;notNull;comment:密码" json:"-"`
	Status        field.Status `gorm:"column:status;notNull;default:0;comment:状态" json:"status"`
	RememberToken string       `gorm:"remember_token;size:500;notNull;default:'';comment:记住密码token" json:"-"`

	Spaces []*Space `gorm:"foreignKey:user_id"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;notNull" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;notNull" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
}
