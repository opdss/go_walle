package model

import (
	"go-walle/app/model/field"
	"time"
)

const (
	TaskStatusWaiting     = 1 //新建提交，等待审核
	TaskStatusAudit       = 2 //审核通过
	TaskStatusReject      = 3 //审核拒绝
	TaskStatusRelease     = 4 //上线发布中
	TaskStatusReleaseFail = 5 //上线失败
	TaskStatusFinish      = 6 //上线完成
)

type Task struct {
	ID            int64 `gorm:"column:id" json:"id"`
	SpaceId       int64 `gorm:"column:space_id;index;not null;comment:所属空间" json:"space_id"`
	ProjectId     int64 `gorm:"column:project_id;not null;comment:所属项目" json:"project_id"`
	UserId        int64 `gorm:"column:user_id;not null;comment:所属用户" json:"user_id"`
	EnvironmentId int64 `gorm:"column:environment_id;not null;comment:所属环境" json:"environment_id"`

	Name        string              `gorm:"column:name;size:100;not null;comment:名称" json:"name"`
	Status      field.Status        `gorm:"column:status;size:1;not null;default:0;comment:状态" json:"status"`
	Version     string              `gorm:"column:version;size:100;not null;default:'';comment:版本号" json:"version"`
	PrevVersion string              `gorm:"column:prev_version;size:100;not null;default:'';comment:上一个版本号" json:"prev_version"`
	ServerIds   field.Slices[int64] `gorm:"column:server_ids" json:"server_ids"`
	CommitId    string              `gorm:"column:commit_id;size:100;not null;default:'';comment:commit哈希" json:"commit_id"`
	Branch      string              `gorm:"column:branch;size:100;not null;default:'';comment:分支" json:"branch"`
	Tag         string              `gorm:"column:tag;size:100;not null;default:'';comment:tag" json:"tag"`
	IsRollback  int                 `gorm:"column:is_rollback" json:"is_rollback"`
	LastError   string              `gorm:"column:last_error" json:"last_error"`
	AuditUserId int64               `gorm:"column:audit_user_id" json:"audit_user_id"`

	Project     Project     `json:"project"`
	User        User        `json:"user"`
	Space       Space       `json:"space"`
	Environment Environment `json:"environment"`
	Servers     []*Server   `gorm:"-" json:"servers"`

	CreatedAt time.Time `gorm:"column:created_at;type:time;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:time;not null" json:"updated_at"`
}
