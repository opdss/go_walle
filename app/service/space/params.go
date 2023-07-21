package space

import (
	"go-walle/app/model/field"
	"go-walle/app/pkg/db"
)

type CreateReq struct {
	Name   string       `json:"name" binding:"required,max=100"`
	Status field.Status `json:"status" binding:"required,status"`
	UserId int64        `json:"user_id" binding:"required"`
}

type UpdateReq struct {
	ID     int64        `json:"id" validate:"required"`
	Name   string       `json:"name" binding:"required,max=100"`
	Status field.Status `json:"status" binding:"required,status"`
	UserId int64        `json:"user_id" binding:"required"`
}

func (r *UpdateReq) Fields() []string {
	return []string{"name", "status", "user_id"}
}

type ListReq struct {
	db.Paginator
}
