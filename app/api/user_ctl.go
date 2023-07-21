package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/model"
	"go-walle/app/service/user"
	"strconv"
)

type UserCtl struct {
	service *user.Service
}

func (ctl *UserCtl) Create(ctx *gin.Context) {
	params := user.CreateReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *UserCtl) List(ctx *gin.Context) {
	params := user.ListReq{}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *UserCtl) Update(ctx *gin.Context) {
	params := user.UpdateReq{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Update(&params), nil)
}

func (ctl *UserCtl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams)
		return
	}
	response.Response(ctx, ctl.service.Delete(int64(n)), nil)
}

func (ctl *UserCtl) Options(ctx *gin.Context) {
	params := user.ListReq{}
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
	res := response.DataOptions{Total: total, Options: slices.Map(items, func(item *model.User, k int) response.DataOption {
		return response.DataOption{
			Text:   fmt.Sprintf("%s(%s)", item.Username, item.Email),
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	response.Success(ctx, res)
}
