package user

import (
	"go-walle/app/model/field"
	"go-walle/app/pkg/db"
)

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Remember bool   `json:"remember" binding:"omitempty"`
}

type LoginRes struct {
	UserId             int64  `json:"user_id"`
	Token              string `json:"token"`
	TokenExpire        int64  `json:"token_expire"`
	RefreshToken       string `json:"refresh_token"`
	RefreshTokenExpire int64  `json:"refresh_token_expire"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token"`
}

type GetUserInfoRes struct {
	UserID         int64        `json:"user_id"`
	Username       string       `json:"username"`
	Email          string       `json:"email"`
	Role           string       `json:"role"`
	Status         field.Status `json:"status"`
	CurrentSpaceId int64        `json:"current_space_id"`
	Spaces         SpaceItems   `json:"spaces"`
}

type CreateReq struct {
	Username string       `json:"username" binding:"required,max=30"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=6,max=50"`
	Status   field.Status `json:"status" binding:"omitempty"`
}

type ListReq struct {
	Keyword string `json:"keyword" query:"keyword" form:"keyword" binding:"omitempty"`
	db.Paginator
}

type MemberListReq struct {
	SpaceId int64
	db.Paginator
}

type MemberListItem struct {
}

type UpdateReq struct {
	ID       int64        `json:"id"  binding:"required,min=1"`
	Username string       `json:"username" binding:"required,max=30"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"omitempty,min=6,max=50"`
	Status   field.Status `json:"status" binding:"omitempty"`
}

type SpaceItem struct {
	SpaceName string       `json:"space_name"`
	SpaceId   int64        `json:"space_id"`
	Status    field.Status `json:"status"`
	Role      string       `json:"role"`
}

type SpaceItems []*SpaceItem

func (s SpaceItems) Default(spaceId int64) *SpaceItem {
	if s == nil || len(s) == 0 {
		return nil
	}
	if spaceId != 0 {
		for _, v := range s {
			if v.SpaceId == spaceId {
				return v
			}
		}
	}
	return s[0]
}
