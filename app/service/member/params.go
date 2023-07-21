package member

import (
	"go-walle/app/model/field"
	"go-walle/app/pkg/db"
	"time"
)

type ListReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`
	db.Paginator
}

type StoreReq struct {
	SpaceId int64  `json:"-" binding:"required,gt=0"`
	UserId  int64  `json:"user_id" binding:"required"`
	Role    string `json:"role" binding:"required"`
}

type ListItem struct {
	SpaceId   int64        `json:"space_id"`
	UserId    int64        `json:"user_id"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Role      string       `json:"role"`
	Status    field.Status `json:"status" `
	CreatedAt time.Time    `json:"created_at" `
}
