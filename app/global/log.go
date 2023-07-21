package global

import (
	"go-walle/app/pkg/log"
	"go.uber.org/zap"
)

var Log *zap.Logger

func initLog(conf *log.Config) (err error) {
	Log = log.NewLog(conf)
	return
}
