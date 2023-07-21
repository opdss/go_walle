package api

import (
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/response"
	"go-walle/app/service/common"
	"go-walle/app/service/user"
)

type LoginCtl struct {
	service *user.Service
}

func (ctl *LoginCtl) Login(ctx *gin.Context) {
	params := user.LoginReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	res, err := ctl.service.Login(&params)
	response.Response(ctx, err, res)
}

func (ctl *LoginCtl) Logout(ctx *gin.Context) {
	response.Response(ctx, ctl.service.Logout(ctx2.UserId(ctx)), nil)
}

func (ctl *LoginCtl) RefreshToken(ctx *gin.Context) {
	params := user.RefreshTokenReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	res, err := ctl.service.RefreshToken(&params)
	response.Response(ctx, err, res)
}

func (ctl *LoginCtl) UserInfo(ctx *gin.Context) {
	spaceAndId := common.SpaceWithId{
		SpaceId: ctx2.GetSpaceId(ctx),
		ID:      ctx2.UserId(ctx),
	}
	res, err := ctl.service.UserInfo(&spaceAndId)
	response.Response(ctx, err, res)
}
