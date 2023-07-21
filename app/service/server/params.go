package server

import "go-walle/app/pkg/db"

type CreateReq struct {
	SpaceId     int64  `json:"-" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required,max=100"`
	User        string `json:"user" binding:"required,max=30"`
	Host        string `json:"host" binding:"required,ip"`
	Port        int    `json:"port" binding:"required,min=22,max=65535"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type UpdateReq struct {
	SpaceId     int64  `json:"-" binding:"required,gt=0"`
	ID          int64  `json:"id" binding:"required,gt=0"`
	Name        string `json:"name" binding:"required,max=100"`
	User        string `json:"user" binding:"required,max=30"`
	Host        string `json:"host" binding:"required,ip"`
	Port        int    `json:"port" binding:"required,min=22,max=65535"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

func (r *UpdateReq) Fields() []string {
	return []string{"name", "user", "host", "port", "description"}
}

type SetAuthorizedReq struct {
	SpaceId  int64  `json:"-" binding:"required,gt=0"`
	ID       int64  `json:"id" binding:"required,gt=0"`
	Password string `json:"password" binding:"required,max=100"`
}

type ListReq struct {
	SpaceId int64 `json:"-" binding:"required,gt=0"`
	//Name string `json:"name" binding:"omitempty,max=100" search:"table:servers;column:name;type:contains"`
	//Host string `json:"host" binding:"omitempty,ip"  search:"table:servers;column:host;type:exact"`
	//common.Order
	db.Paginator
}
