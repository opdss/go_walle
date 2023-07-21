package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/model"
	"go-walle/app/service/environment"
)

type EnvironmentCtl struct {
	service *environment.Service
}

func (ctl *EnvironmentCtl) Create(ctx *gin.Context) {
	params := environment.CreateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *EnvironmentCtl) List(ctx *gin.Context) {
	params := environment.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *EnvironmentCtl) Update(ctx *gin.Context) {
	params := environment.UpdateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Update(&params), nil)
}

func (ctl *EnvironmentCtl) Delete(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Delete(spaceAndId), nil)
}

func (ctl *EnvironmentCtl) Options(ctx *gin.Context) {
	params := environment.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	res := response.DataOptions{Total: total, Options: slices.Map(items, func(item *model.Environment, k int) response.DataOption {
		return response.DataOption{
			Text:   item.Name,
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	response.Success(ctx, res)
}
