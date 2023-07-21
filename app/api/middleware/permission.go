package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	ctx2 "go-walle/app/api/ctx"
	"go-walle/app/internal/constants"
	"go-walle/app/service/user"
	"log"
)

func Permission(userService *user.Service, role constants.Role) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Println("middleware Permission start")
		userId := ctx2.UserId(ctx)
		spaceId := ctx2.GetSpaceId(ctx)
		if !constants.IsSuperUser(userId) {
			if spaceId == 0 {
				_ = ctx.AbortWithError(400, errors.New("未选择空间"))
				return
			}
			member, err := userService.SpaceById(userId, spaceId)
			if err != nil {
				_ = ctx.AbortWithError(400, errors.New("空间选择错误"))
				return
			}
			currRole := member.Role
			if constants.Role(currRole).Level() < role.Level() {
				_ = ctx.AbortWithError(401, errors.New("你没有权限访问该空间，请联系相关负责人"))
				return
			}
		}
		//ctx2.SetRole(ctx, currRole)
		ctx.Next()
		log.Println("middleware Permission end")
	}
}
