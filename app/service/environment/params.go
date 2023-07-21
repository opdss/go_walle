package environment

import (
	"go-walle/app/model/field"
	"go-walle/app/pkg/db"
)

type CreateReq struct {
	SpaceId     int64        `json:"-" binding:"required,gt=0"`
	Name        string       `json:"name" binding:"required,max=100"`
	Status      field.Status `json:"status" binding:"required,status"`
	Description string       `json:"description" binding:"omitempty,max=500"`
	Color       string       `json:"color" binding:"omitempty,rgb"`
}

type UpdateReq struct {
	SpaceId     int64        `json:"-" binding:"required,gt=0"`
	ID          int64        `json:"id" binding:"required"`
	Name        string       `json:"name" binding:"required,max=100"`
	Status      field.Status `json:"status" binding:"required,status"`
	Description string       `json:"description" binding:"omitempty,max=500"`
	Color       string       `json:"color" binding:"omitempty,rgb"`
}

func (r *UpdateReq) Fields() []string {
	return []string{"name", "status", "description", "color"}
}

type ListReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`
	db.Paginator
}
