package model

import (
	"go-walle/app/model/field"
	"gorm.io/gorm"
	"time"
)

type Environment struct {
	ID          int64        `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	SpaceId     int64        `gorm:"column:space_id;index;not null;comment:所属空间" json:"space_id"`
	Name        string       `gorm:"column:name;size:100;not null;comment:名称" json:"name"`
	Status      field.Status `gorm:"column:status;size:1;not null;default:0;comment:状态" json:"status"`
	Description string       `gorm:"column:description;size:500;not null;default:'';comment:简介说明" json:"description"`
	Color       string       `gorm:"column:color;size:10;not null;default:'';comment:主题色" json:"color"`

	Space    Space      `json:"space"`
	Projects []*Project `json:"projects"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
