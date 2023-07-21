package services

import (
	"context"
	"go-walle/app/global"
	"go-walle/app/model"
	"go-walle/app/pkg/ssh"
	ssh2 "golang.org/x/crypto/ssh"
	"os/exec"
	"time"
)

type Record struct {
	model  *model.Record
	server *model.Server
	envs   *ssh.Envs
}

func NewRecordLocal(taskId, userId int64, cmd string, envs *ssh.Envs) *Record {
	return &Record{
		model: &model.Record{
			UserId:   userId,
			TaskId:   taskId,
			ServerId: 0,
			Envs:     envs.SliceKV(),
		},
	}
}

func NewRecordRemote(taskId, userId int64, cmd string, server *model.Server, envs *ssh.Envs) *Record {
	return &Record{}
}

func (r *Record) Run(ctx context.Context) (err error) {
	startT := time.Now()
	var command ssh.Command
	if r.server == nil {
		command = ssh.NewLocalExec()
	} else {
		command, err = global.Ssh.NewRemoteExec(ssh.ServerConfig{
			Host:     r.server.Host,
			User:     r.server.User,
			Password: "",
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
