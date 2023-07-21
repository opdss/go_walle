package ssh

import (
	"context"
	"fmt"
	"strings"
)

type RemoteExec struct {
	client    *client
	envs      *Envs
	ctx       context.Context
	closeChan chan int
}

func (e *RemoteExec) Close() error {
	e.client.done()
	return nil
}

func (e *RemoteExec) WithCtx(ctx context.Context) Command {
	e.ctx = ctx
	return e
}
func (e *RemoteExec) WithEnvs(envs *Envs) Command {
	e.envs = envs
	return e
}

func (e *RemoteExec) Run(cmd string) ([]byte, error) {
	sess, err := e.client.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = sess.Close()
	}()
	if !e.envs.Empty() {
		cmd = fmt.Sprintf("%s && %s", strings.Join(e.envs.SliceKV(), " "), cmd)
	}
	if e.ctx != nil {
		e.closeChan = make(chan int, 1)
		defer close(e.closeChan)
		go func() {
			select {
			case _, ok := <-e.closeChan:
				if !ok {
					return
				}
			case <-e.ctx.Done():
				_ = sess.Close()
				return
			}
		}()
	}
	return sess.CombinedOutput(cmd)
}
