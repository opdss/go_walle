package command

import (
	"context"
	"os/exec"
)

type Local struct {
	cmd  string
	envs *Envs
}

func NewLocal(cmd string, env *Envs) {

}

func (l *Local) Run(ctx context.Context) error {
	command := exec.CommandContext(ctx, "bash", "-c", l.cmd)
	command.Env = l.envs.SliceKV()
	return command.Run()
}

func (l *Local) CombinedOutput(ctx context.Context) ([]byte, error) {
	command := exec.CommandContext(ctx, "bash", "-c", l.cmd)
	command.Env = l.envs.SliceKV()
	return command.CombinedOutput()
}
