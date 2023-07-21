package api

import (
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/go-helper/slices"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/errcode"
	"go-walle/app/internal/response"
	"go-walle/app/model"
	"go-walle/app/service/project"
)

type ProjectCtl struct {
	service *project.Service
}

func (ctl *ProjectCtl) Create(ctx *gin.Context) {
	params := project.CreateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Create(&params), nil)
}

func (ctl *ProjectCtl) List(ctx *gin.Context) {
	params := project.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBind(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	total, items, err := ctl.service.List(&params)
	response.PageData(ctx, total, items, err)
}

func (ctl *ProjectCtl) Update(ctx *gin.Context) {
	params := project.UpdateReq{SpaceId: ctx2.GetSpaceId(ctx)}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Update(&params), nil)
}

func (ctl *ProjectCtl) Delete(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	response.Response(ctx, ctl.service.Delete(spaceAndId), nil)
}

func (ctl *ProjectCtl) Detail(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := ctl.service.Detail(spaceAndId)
	response.Response(ctx, err, res)
}

func (ctl *ProjectCtl) Options(ctx *gin.Context) {
	params := project.ListReq{SpaceId: ctx2.GetSpaceId(ctx)}
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
	res := response.DataOptions{Total: total, Options: slices.Map(items, func(item *model.Project, k int) response.DataOption {
		return response.DataOption{
			Text:   item.Name,
			Value:  item.ID,
			Status: item.Status,
		}
	})}
	response.Success(ctx, res)
}

// Detection 检测项目
func (ctl *ProjectCtl) Detection(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := ctl.service.Detection(spaceAndId)
	response.Response(ctx, err, res)
}

// Branches 分支列表
func (ctl *ProjectCtl) Branches(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := ctl.service.GetBranches(spaceAndId)
	response.Response(ctx, err, res)
}

// Tags tags列表
func (ctl *ProjectCtl) Tags(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	if err != nil {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := ctl.service.GetTags(spaceAndId)
	response.Response(ctx, err, res)
}

// Commits 提交记录
func (ctl *ProjectCtl) Commits(ctx *gin.Context) {
	spaceAndId, err := ctx2.GetSpaceWithId(ctx)
	branch := ctx.Query("branch")
	if err != nil || branch == "" {
		response.Fail(ctx, errcode.ErrInvalidParams.Wrap(err))
		return
	}
	res, err := ctl.service.GetCommits(spaceAndId, ctx.Query("branch"))
	response.Response(ctx, err, res)
}
