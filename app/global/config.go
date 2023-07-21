package global

import (
	errs2 "github.com/zeebo/errs"
	"go-walle/app/pkg/db"
	"go-walle/app/pkg/jwt"
	"go-walle/app/pkg/log"
	"go-walle/app/pkg/repo"
	"go-walle/app/pkg/ssh"
)

var Cfg *Config

type Config struct {
	Api struct {
		Address string `help:"监听地址" devDefault:"0.0.0.0:8989" default:"0.0.0.0:8080"`
	}
	Db   db.Config
	Repo repo.Config
	JWT  jwt.Config
	Log  log.Config
	Ssh  ssh.Config
}

func (c *Config) Init() {
	Cfg = c
	errs := errs2.Group{}
	errs.Add(
		initLog(&c.Log),
		initDB(&c.Db),
		initJwt(&c.JWT),
		initRepo(&c.Repo),
		initSsh(&c.Ssh),
	)
	if errs.Err() != nil {
		panic(errs.Err())
	}
}
