package model

import (
	"time"
)

type Member struct {
	ID      int64  `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	UserId  int64  `gorm:"column:user_id;not null;uniqueIndex:space_user;comment:用户" json:"user_id"`
	SpaceId int64  `gorm:"column:space_id;not null;uniqueIndex:space_user;comment:空间" json:"space_id"`
	Role    string `gorm:"column:role;size:20;not null;comment:角色" json:"role"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`

	Space Space
	User  User
}
