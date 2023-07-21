package deploy

import (
	"go-walle/app/model"
	"go-walle/app/pkg/db"
)

type CreateReq struct {
	UserId      int64   `json:"-" binding:"required,gt=0"`
	SpaceId     int64   `json:"-" binding:"required,gt=0"`
	ProjectId   int64   `json:"project_id" binding:"required,gt=0"`
	Name        string  `json:"name" binding:"required,max=100"`
	Tag         string  `json:"tag" binding:"omitempty,max=50"`
	Branch      string  `json:"branch" binding:"omitempty,max=50"`
	CommitId    string  `json:"commit_id" binding:"omitempty,max=50"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	ServerIds   []int64 `json:"server_ids" binding:"omitempty"`
}

type ListReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`
	db.Paginator
}

type AuditReq struct {
	SpaceId     int64 `json:"-" binding:"required,gt=0"`
	AuditUserId int64 `json:"-" binding:"required,gt=0"`
	ID          int64 `json:"-" binding:"required,gt=0"`
	Audit       bool  `json:"audit" `
}

const TaskConsoleMsgRecords = "records"
const TaskConsoleMsgRecord = "record"
const TaskConsoleMsgAppend = "append"

type ConsoleMsg struct {
}

type TaskConsoleMsg struct {
	Type    string          `json:"type"`
	Records []*model.Record `json:"records,omitempty"`
}
