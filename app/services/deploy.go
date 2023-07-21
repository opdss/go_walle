package services

import (
	"github.com/zeebo/errs"
	"golang.org/x/crypto/ssh"
	"sync"
)

var Error = errs.Class("Deploy")

type Deploy struct {
	mux        *sync.RWMutex
	tasks      map[int64]*task
	sshClients []*ssh.Client
}

// Start 开始部署
func (d *Deploy) Start(taskId int64) error {
	d.mux.Lock()
	defer d.mux.Unlock()
	if _, ok := d.tasks[taskId]; ok {
		return Error.New("task[%d]已经开始部署", taskId)
	}
	d.tasks[taskId] = &task{}
	return d.tasks[taskId].Start()
}

// Stop 中止部署
func (d *Deploy) Stop(taskId int64) {
	d.mux.Lock()
	defer d.mux.Unlock()
	if v, ok := d.tasks[taskId]; ok {
		v.Stop()
		delete(d.tasks, taskId)
		return
	}
}

// Output 获取输出
func (d *Deploy) Output(taskId int64) (<-chan []byte, error) {
	c := make(chan []byte, 0)
	return c, nil
}
