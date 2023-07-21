package deploy

import (
	"fmt"
	"go-walle/app/model"
	"go-walle/app/pkg/db"
	"testing"
)

var _db *db.DB

func init() {
	_db, _ = db.NewGormDB(&db.Config{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "agEaB&^23fB",
		Database: "go-walle",
	}, nil)
}

func TestTask(t *testing.T) {
	var taskModel model.Task
	err := _db.Where("id = 45").First(&taskModel).Error
	fmt.Println(taskModel.ID)
	if err != nil {
		t.Error("查询数据错误", err)
	}
	task, err := CreateDeployTask(&taskModel, 1)
	if err != nil {
		t.Error("创建任务错误", err)
	}
	err = task.Start()
	if err != nil {
		t.Error("启动任务失败", err)
	}
}
