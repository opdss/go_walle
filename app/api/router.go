package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	"go-walle/app/api/middleware"
	"go-walle/app/global"
	"go-walle/app/internal/constants"
	"go-walle/app/service/deploy"
	"go-walle/app/service/environment"
	"go-walle/app/service/member"
	"go-walle/app/service/project"
	server2 "go-walle/app/service/server"
	"go-walle/app/service/space"
	"go-walle/app/service/user"
	"io"
	"mime"
	"path"
	"strings"
)

func ApiRoutes(e *gin.Engine, server *Server) {
	//这三个是静态文件的路由
	e.GET("/", func(ctx *gin.Context) {
		fileHandle(ctx, server.rootFs, "index.html")
	})
	e.GET("/:file", func(ctx *gin.Context) {
		fileHandle(ctx, server.rootFs, ctx.Param("file"))
	})
	e.GET("/:file/*child", func(ctx *gin.Context) {
		fileHandle(ctx, server.rootFs, ctx.Param("file")+ctx.Param("child"))
	})
	e.GET("/assets/*file", func(ctx *gin.Context) {
		fileHandle(ctx, server.assetsFs, "assets"+ctx.Param("file"))
	})

	r := e.Group("/api")
	apiRoutes(r, server)
}

// 静态文件处理
func fileHandle(ctx *gin.Context, fs *embed.FS, file string) {
	file = "web/dist/" + strings.TrimPrefix(file, "/")
	f, err := fs.Open(file)
	if err != nil {
		ctx.AbortWithStatus(404)
		return
	}
	ctx.Header("Content-Type", mime.TypeByExtension(path.Ext(file)))
	_, err = io.Copy(ctx.Writer, f)
	if err != nil {
		_ = ctx.AbortWithError(500, err)
		return
	}
}

func apiRoutes(r *gin.RouterGroup, s *Server) {

	userService := user.NewService(global.Log, global.DB, global.Jwt)

	commonCtl := CommonCtl{}
	r.GET("/version", commonCtl.Version)
	r.GET("/server_info", commonCtl.ServiceInfo)

	loginCtl := LoginCtl{service: userService}
	r.POST("/login", loginCtl.Login)
	r.POST("/refresh_token", loginCtl.RefreshToken)

	authRouter := r.Group("", middleware.Auth)
	authRouter.POST("/logout", loginCtl.Logout)
	authRouter.GET("/user_info", loginCtl.UserInfo)

	superPermRouter := authRouter.Group("", middleware.Permission(userService, constants.RoleSuper))
	ownerPermRouter := authRouter.Group("", middleware.Permission(userService, constants.RoleOwner))
	masterPermRouter := authRouter.Group("", middleware.Permission(userService, constants.RoleMaster))
	//developerPermMid := middleware.Permission(constants.RoleDeveloper)

	//用户管理
	{
		ctl := &UserCtl{service: user.NewService(global.Log, global.DB, global.Jwt)}
		superPermRouter.GET("/user", ctl.List)
		superPermRouter.POST("/user", ctl.Create)
		superPermRouter.DELETE("/user/:id", ctl.Delete)
		superPermRouter.PUT("/user", ctl.Update)
		superPermRouter.GET("/user/options", ctl.Options)
	}

	//成员管理
	{
		ctl := &MemberCtl{service: member.NewService()}
		ownerPermRouter.GET("/member", ctl.List)
		ownerPermRouter.POST("/member", ctl.Store)
		ownerPermRouter.DELETE("/member/:id", ctl.Delete)
	}

	//空间管理, super访问权限
	{
		ctl := &SpaceCtl{service: space.NewService(global.DB)}
		superPermRouter.GET("/space", ctl.List)
		superPermRouter.POST("/space", ctl.Create)
		superPermRouter.DELETE("/space/:id", ctl.Delete)
		superPermRouter.PUT("/space", ctl.Update)
	}

	//服务器管理
	{
		ctl := &ServerCtl{service: server2.NewService(global.Log, global.DB, global.Ssh)}
		ownerPermRouter.GET("/server", ctl.List)
		ownerPermRouter.POST("/server", ctl.Create)
		ownerPermRouter.DELETE("/server/:id", ctl.Delete)
		ownerPermRouter.PUT("/server", ctl.Update)
		//校验连接
		ownerPermRouter.POST("/server/:id/check", ctl.Check)
		//设置免登陆
		ownerPermRouter.POST("/server/set_authorized", ctl.SetAuthorized)
		//websocket 连接终端
		ownerPermRouter.GET("/server/:id/terminal", ctl.Terminal)
	}

	//环境管理
	{
		ctl := &EnvironmentCtl{service: environment.NewService(global.DB)}
		masterPermRouter.GET("/environment", ctl.List)
		masterPermRouter.POST("/environment", ctl.Create)
		masterPermRouter.DELETE("/environment/:id", ctl.Delete)
		masterPermRouter.PUT("/environment", ctl.Update)
		masterPermRouter.GET("/environment/options", ctl.Options)
	}

	//项目管理
	{
		ctl := &ProjectCtl{service: project.NewService(global.DB, global.Ssh, global.Repo)}
		masterPermRouter.GET("/project", ctl.List)
		masterPermRouter.POST("/project", ctl.Create)
		masterPermRouter.DELETE("/project/:id", ctl.Delete)
		masterPermRouter.GET("/project/:id", ctl.Detail)
		masterPermRouter.PUT("/project", ctl.Update)
		masterPermRouter.GET("/project/options", ctl.Options)
		//项目检测
		masterPermRouter.GET("/project/:id/detection", ctl.Detection)
		masterPermRouter.GET("/project/:id/branches", ctl.Branches)
		masterPermRouter.GET("/project/:id/tags", ctl.Tags)
		masterPermRouter.GET("/project/:id/commits", ctl.Commits)
	}

	//部署管理
	{
		ctl := &DeployCtl{service: deploy.NewService()}
		masterPermRouter.GET("/deploy", ctl.List)
		masterPermRouter.GET("/deploy/:id", ctl.Detail)
		masterPermRouter.POST("/deploy", ctl.Create)
		//审核
		masterPermRouter.POST("/deploy/:id/audit", ctl.Audit)
		//发布
		masterPermRouter.GET("/deploy/:id/release", ctl.Release)
		//发布
		masterPermRouter.GET("/deploy/:id/stop_release", ctl.StopRelease)
		//发布
		masterPermRouter.GET("/deploy/:id/rollback", ctl.Rollback)
		//websocket, 部署日志, 将整个部署过程日志输出
		masterPermRouter.GET("/deploy/:id/console", ctl.Console)
	}

}
