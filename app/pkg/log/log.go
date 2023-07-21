package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/wuzfei/go-helper/time"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	levelMap = map[string]zapcore.Level{
		"debug":   zapcore.DebugLevel,
		"error":   zapcore.ErrorLevel,
		"warning": zapcore.WarnLevel,
		"info":    zapcore.InfoLevel,
		"panic":   zapcore.PanicLevel,
		"fatal":   zapcore.FatalLevel,
	}
)

type Config struct {
	File        string `help:"日志输出文件" devDefault:"$ROOT/go_walle.log" default:"$HOME/go_walle.log"`
	FileSize    int    `help:"日志文件大小限制,单位MB" default:"500"`
	FileBackups int    `help:"最大保留日志文件数量" default:"10"`
	FileAge     int    `help:"日志文件保留天数" devDefault:"0" default:"30"`
	Level       string `help:"日志输出级别,可选[debug|info|warning|error|panic|fatal]" default:"debug"`
	Output      string `help:"日志输出方式,可选[any|file|console]"  devDefault:"any" default:"file"`
	Encoder     string `help:"日志输出格式,可选[json|console]" devDefault:"console" default:"json"`
}

func NewLog(cfg *Config) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DefaultDateTimeLayout)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	var encoder zapcore.Encoder
	if cfg.Encoder == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	//文件writeSyncer
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.File,        //日志文件存放目录
		MaxSize:    cfg.FileSize,    //文件大小限制,单位MB
		MaxBackups: cfg.FileBackups, //最大保留日志文件数量
		MaxAge:     cfg.FileAge,     //日志文件保留天数
		Compress:   false,           //是否压缩处理
	})

	level := zapcore.InfoLevel
	if _level, ok := levelMap[cfg.Level]; ok {
		level = _level
	}
	consoleOutput := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	fileOutput := zapcore.NewCore(encoder, fileWriteSyncer, level)
	var core zapcore.Core
	switch cfg.Output {
	case "any":
		core = zapcore.NewTee(consoleOutput, fileOutput)
		break
	case "file":
		core = fileOutput
		break
	default:
		core = consoleOutput
	}
	return zap.New(core)
}
