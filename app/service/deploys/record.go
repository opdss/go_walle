package deploys

import (
	"context"
	"encoding/json"
	"go-walle/app/global"
	"go-walle/app/model"
	"go-walle/app/pkg/command"
	"go-walle/app/pkg/ssh"
	"go.uber.org/zap"
	ssh2 "golang.org/x/crypto/ssh"
	"os/exec"
	"time"
)

type record struct {
	model  *model.Record
	server *model.Server
	envs   *ssh.Envs
	cmd    command.Command
}

func NewRecord(typ int, taskId, userId int64, cmd string, server *model.Server, envs *ssh.Envs) *record {
	if envs == nil {
		envs = ssh.NewEnvs()
	}
	r := &record{
		model: &model.Record{
			Type:    typ,
			UserId:  userId,
			Status:  -1,
			Command: cmd,
			Envs:    envs.SliceKV(),
			TaskId:  taskId,
		},
		server: server,
		envs:   envs,
	}
	if server != nil {
		r.model.ServerId = r.server.ID
	}
	return r
}

func (r *record) Run(ctx context.Context) (err error) {
	startT := time.Now()
	var command ssh.Command
	if r.server == nil {
		command = ssh.NewLocalExec()
	} else {
		command, err = global.Ssh.NewRemoteExec(ssh.ServerConfig{
			Host:     r.server.Host,
			User:     r.server.User,
			Password: "ipfs",
			Port:     r.server.Port,
		})
	}
	if err == nil {
		var output []byte
		output, err = command.WithEnvs(r.envs).WithCtx(ctx).Run(r.model.Command)
		r.model.Output = string(output)
	}
	if err != nil {
		if e, ok := err.(*ssh2.ExitError); ok {
			r.model.Status = e.ExitStatus()
		} else if e, ok := err.(*exec.ExitError); ok {
			r.model.Status = e.ExitCode()
		} else {
			r.model.Status = 255
		}
	} else {
		r.model.Status = 0
	}
	r.model.RunTime = time.Since(startT).Milliseconds()
	_ = r.save()
	return err
}

func (r *record) Save(status int, output *string, runtime int64) error {
	r.model.RunTime = runtime
	r.model.Status = status
	r.model.Output = *output
	return r.save()
}

func (r *record) Output() string {
	return r.model.Output
}

func (r *record) save() error {
	err := global.DB.Create(r.model).Error
	if err != nil {
		obj, _ := json.Marshal(r.model)
		global.Log.Error("保存执行记录失败", zap.ByteString("record", obj))
	}
	return err
}
