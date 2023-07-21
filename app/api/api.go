package api

import (
	"context"
	"embed"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wuzfei/cfgstruct/cfgstruct"
	"go-walle/app/global"
	"go-walle/app/internal/validate"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
)

type Server struct {
	config   *global.Config
	rootFs   *embed.FS
	assetsFs *embed.FS
	server   http.Server
}

func NewServer(conf *global.Config, rootfs *embed.FS, assets *embed.FS) *Server {
	server := &Server{
		rootFs:   rootfs,
		assetsFs: assets,
		config:   conf,
	}
	return server
}

func (s *Server) Run(ctx context.Context) error {
	if cfgstruct.DefaultsType() == cfgstruct.DefaultsRelease {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	ApiRoutes(engine, s)
	// 注册自定义验证标签
	if err := validate.RegisterValidation(); err != nil {
		return err
	}
	s.server.Handler = engine

	listener, err := net.Listen("tcp", s.config.Api.Address)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return s.server.Shutdown(context.Background())
	})
	group.Go(func() error {
		defer cancel()
		_err := s.server.Serve(listener)
		if errors.Is(_err, http.ErrServerClosed) {
			_err = nil
		}
		return _err
	})
	return group.Wait()
}
