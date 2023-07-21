package api

import (
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/service/deploy"
	"strconv"
)

type DeployCtl struct {
	service *deploy.Service
}

func (ctl *DeployCtl) Create(ctx *gin.Context) {
	params := deploy.CreateReq{SpaceId: ctx2.GetSpaceId(ctx), UserId: ctx2.UserId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *DeployCtl) List(ctx *gin.Context) {
	params := deploy.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *DeployCtl) Detail(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	data, err := ctl.service.Detail(spaceAndId)
	response.Response(ctx, err, data)
}

func (ctl *DeployCtl) Audit(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	params := deploy.AuditReq{SpaceId: ctx2.GetSpaceId(ctx), AuditUserId: ctx2.UserId(ctx), ID: int64(n)}
	err = ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Audit(&params), nil)
}

func (ctl *DeployCtl) Release(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = ctl.service.Release(spaceAndId, ctx2.UserId(ctx))
	response.Response(ctx, err, nil)
}

func (ctl *DeployCtl) StopRelease(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = ctl.service.Release(spaceAndId, ctx2.UserId(ctx))
	response.Response(ctx, err, nil)
}

func (ctl *DeployCtl) Rollback(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	err = ctl.service.Release(spaceAndId, ctx2.UserId(ctx))
	response.Response(ctx, err, nil)
}

func (ctl *DeployCtl) Console(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	wsConn, err := ctx2.UpGrader(ctx)
	if err != nil {
		response.Fail(ctx, err)
		return
	}
	defer func() {
		_ = wsConn.Close()
	}()
	ctl.service.Console(wsConn, spaceAndId)
}
