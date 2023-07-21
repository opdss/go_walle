package api

import (
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/service/member"
)

type MemberCtl struct {
	service *member.Service
}

func (ctl *MemberCtl) Store(ctx *gin.Context) {
	params := member.StoreReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Store(params), nil)
}

func (ctl *MemberCtl) List(ctx *gin.Context) {
	params := member.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(params)
	response.PageData(ctx, total, items, err)
}

func (ctl *MemberCtl) Delete(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Delete(spaceAndId), nil)
}
