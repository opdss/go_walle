package deploy

import (
	"errors"
	"go-walle/app/model"
	"sync"
)

var (
	mux         *sync.Mutex
	deployTasks map[int64]*Task
)

func init() {
	mux = &sync.Mutex{}
	deployTasks = map[int64]*Task{}
}

func CreateDeployTask(model *model.Task, userId int64) (*Task, error) {
	mux.Lock()
	defer mux.Unlock()
	if _, ok := deployTasks[model.ID]; ok {
		return nil, errors.New("该任务已经在部署中")
	}
	o := NewTask(model, userId)
	deployTasks[model.ID] = o
	return o, nil
}

func GetDeployTask(taskId int64) *Task {
	mux.Lock()
	defer mux.Unlock()
	if o, ok := deployTasks[taskId]; ok {
		return o
	}
	return nil
}

func removeDeployTask(taskId int64) {
	mux.Lock()
	defer mux.Unlock()
	if _, ok := deployTasks[taskId]; ok {
		delete(deployTasks, taskId)
	}
}
