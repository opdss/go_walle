package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/global"
	"go-walle/app/internal/errcode"
	"go-walle/app/model/field"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

type DataOption struct {
	Text   string       `json:"text"`   //显示文本
	Value  interface{}  `json:"value"`  //数据id
	Status field.Status `json:"status"` //是否可选 ，1可选 2不可选
	Other  interface{}  `json:"other"`  //其他
}

type DataOptions struct {
	Total   int64        `json:"total"`   //总数量
	Options []DataOption `json:"options"` //选项数据
}

type PageResult struct {
	Total int64 `json:"total"`
	Items any   `json:"items"`
}

func Success(ctx *gin.Context, data any) {
	ctx.AbortWithStatusJSON(200, Data{Code: 0, Message: "success", Result: data})
}

func PageData(ctx *gin.Context, total int64, items any, err error) {
	if err != nil {
		Fail(ctx, err)
		return
	}
	ctx.AbortWithStatusJSON(200, Data{Code: 0, Message: "success", Result: PageResult{
		Total: total,
		Items: items,
	}})
}

func SuccessPageData(ctx *gin.Context, total int64, items any) {
	ctx.AbortWithStatusJSON(200, Data{Code: 0, Message: "success", Result: PageResult{
		Total: total,
		Items: items,
	}})
}

func Fail(ctx *gin.Context, err error, data ...any) {
	code := -1
	msg := err.Error()
	if e, ok := err.(errcode.ErrCode); ok {
		code = int(e)
	} else if e, ok := err.(errcode.ErrWrap); ok {
		code = e.ErrCode()
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		msg = "数据无权限或者不存在"
	}
	global.Log.Error("Response Error",
		zap.Int64("UserId", ctx2.UserId(ctx)),
		zap.Int64("SpaceId", ctx2.GetSpaceId(ctx)),
		zap.Int("ErrCode", code),
		zap.Error(err))
	ctx.AbortWithStatusJSON(200, Data{Code: code, Message: msg, Result: data})
}

func Response(ctx *gin.Context, err error, data any) {
	if err != nil {
		Fail(ctx, err)
	} else {
		Success(ctx, data)
	}
}
