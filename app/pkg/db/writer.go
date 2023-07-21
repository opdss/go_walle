package db

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"time"
)

type writer struct {
	*zap.Logger
}

func (l *writer) Printf(fm string, args ...interface{}) {
	l.Info(fmt.Sprintf(fm, args...))
}

func getLogInterface(zapLog *zap.Logger, logLevel string) logger.Interface {
	level := logger.Silent
	switch logLevel {
	case "error":
		level = logger.Error
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	}
	if zapLog == nil {
		return logger.Default.LogMode(level)
	} else {
		return logger.New(&writer{zapLog.Named("DB")}, logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		})
	}
}
