package api

import (
	"github.com/gin-gonic/gin"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/service/space"
	"strconv"
)

type SpaceCtl struct {
	service *space.Service
}

func (ctl *SpaceCtl) Create(ctx *gin.Context) {
	params := space.CreateReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *SpaceCtl) List(ctx *gin.Context) {
	params := space.ListReq{}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *SpaceCtl) Update(ctx *gin.Context) {
	params := space.UpdateReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Update(&params), nil)
}

func (ctl *SpaceCtl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
	}
	response.Response(ctx, ctl.service.Delete(int64(n)), nil)
}
