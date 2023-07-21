package model

import (
	"go-walle/app/model/field"
	"time"
)

const (
	RecordTypeDefault = iota
	RecordTypePrevDeploy
	RecordTypeDeploy
	RecordTypePostDeploy
	RecordTypePrevRelease
	RecordTypeRelease
	RecordTypePostRelease
)

type Record struct {
	ID       int64                `gorm:"column:id;primaryKey;autoIncrement;" json:"id"`
	Type     int                  `gorm:"column:type" json:"type"`
	UserId   int64                `gorm:"column:user_id" json:"user_id"`
	ServerId int64                `gorm:"column:server_id" json:"server_id"`
	TaskId   int64                `gorm:"column:task_id" json:"task_id"`
	Envs     field.Slices[string] `gorm:"column:envs" json:"envs"`
	RunTime  int64                `gorm:"column:run_time" json:"run_time"`
	Status   int                  `gorm:"column:status" json:"status"`
	Command  string               `gorm:"column:command" json:"command"`
	Output   string               `gorm:"column:output" json:"output"`

	Server Server `json:"server"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`
}
