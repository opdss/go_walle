package global

import (
	"go-walle/app/pkg/db"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initDB(conf *db.Config) (err error) {
	DB, err = db.NewGormDB(conf, Log)
	return
}
