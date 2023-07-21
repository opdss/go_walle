package model

import (
	"fmt"
	"go-walle/app/model/field"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	ID          int64        `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId     int64        `gorm:"column:space_id;index;not null;comment:所属空间" json:"space_id"`
	Name        string       `gorm:"column:name;size:100;not null;comment:名称" json:"name"`
	User        string       `gorm:"column:user;size:100;not null;comment:用户名" json:"user"`
	Host        string       `gorm:"column:host;size:100;not null;comment:主机" json:"host"`
	Port        int          `gorm:"column:port;size:4;not null;default:22;comment:端口" json:"port"`
	Status      field.Status `gorm:"column:status;size:1;not null;default:0;comment:状态" json:"status"`
	Description string       `gorm:"column:description;size:500;not null;default:'';comment:简介说明" json:"description"`

	Projects []*Project `gorm:"many2many:project_server" json:"projects"`
	Tasks    []*Task    `gorm:"many2many:task_server" json:"tasks"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (receiver *Server) Hostname() string {
	return fmt.Sprintf("%s@%s", receiver.User, receiver.Host)
}
