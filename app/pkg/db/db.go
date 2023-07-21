package db

import (
	"fmt"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Mysql = "mysql"
const Postgresql = "postgres"
const Sqlite = "sqlite"

var ErrDB = errs.Class("DB")

type Config struct {
	Driver   string `help:"数据库驱动" default:"mysql"`
	Host     string `help:"数据库地址" devDefault:"192.168.43.26" default:"localhost"`
	Port     int    `help:"数据库端口" devDefault:"3307" default:"3306"`
	Username string `help:"数据库帐号" default:"root"`
	Password string `help:"数据库密码" default:"root"`
	Database string `help:"数据库名称" default:"test"`
	Charset  string `help:"数据库编码" default:"utf8mb4"`
	SslMode  string `help:"pg用" default:"false"`
	TimeZone string `help:"时区" default:"Asia/Shanghai"`
	File     string `help:"数据库，sqlite用" default:"$ROOT/go_walle.db"`
	LogLevel string `help:"数据库日志打印级别,默认为空,可选[error|warn|info]" devDefault:"" default:"warn"`
}

func (c *Config) GetDsn() (dsn string, err error) {
	switch c.Driver {
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t",
			c.Username, c.Password, c.Host, c.Port, c.Database, c.Charset, true)
	case Postgresql:
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			c.Host, c.Username, c.Password, c.Database, c.Port, c.SslMode, c.TimeZone)
	case Sqlite:
		dsn = c.File
	default:
		err = ErrDB.New("数据库驱动错误：%s", c.Driver)
	}
	return
}

func (c *Config) Dialector() (dial gorm.Dialector, err error) {
	var dsn string
	dsn, err = c.GetDsn()
	if err != nil {
		return
	}
	switch c.Driver {
	case Mysql:
		dial = mysql.New(mysql.Config{
			DSN:                       dsn,
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		})
		return
	case Postgresql:
		dial = postgres.New(postgres.Config{
			DSN: dsn,
		})
		return
	case Sqlite:
		dial = sqlite.Open(dsn)
		return
	}
	return nil, ErrDB.New("数据库驱动错误：%s", c.Driver)
}
func NewGormDB(cfg *Config, zapLog *zap.Logger) (*gorm.DB, error) {
	dail, err := cfg.Dialector()
	if err != nil {
		return nil, ErrDB.Wrap(err)
	}
	//todo logger
	return gorm.Open(dail, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   getLogInterface(zapLog, cfg.LogLevel),
	})
}
