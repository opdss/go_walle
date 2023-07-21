package api

import (
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/global"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/service/server"
	"go.uber.org/zap"
)

type ServerCtl struct {
	service *server.Service
}

func (ctl *ServerCtl) Create(ctx *gin.Context) {
	params := server.CreateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *ServerCtl) List(ctx *gin.Context) {
	params := server.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *ServerCtl) Update(ctx *gin.Context) {
	params := server.UpdateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Update(&params), nil)
}

func (ctl *ServerCtl) Delete(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Delete(spaceAndId), nil)
}

func (ctl *ServerCtl) Check(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Check(spaceAndId), nil)
}

func (ctl *ServerCtl) SetAuthorized(ctx *gin.Context) {
	params := server.SetAuthorizedReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	response.Response(ctx, ctl.service.SetAuthorized(&params), nil)
}

func (ctl *ServerCtl) Terminal(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	wsConn, err := ctx2.UpGrader(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrServer.Wrap(err))
	}
	defer func() {
		_ = wsConn.Close()
	}()
	if err = ctl.service.Terminal(wsConn, spaceAndId, ctx2.Username(ctx)); err != nil {
		global.Log.Error("terminal error", zap.Error(err))
	}
}
