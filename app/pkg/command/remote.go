package command

import (
	"context"
)

type Server struct {
}

type Remote struct {
	cmd    string
	envs   *Envs
	server *Server
}

func (r *Remote) Run(ctx context.Context) error {
	return nil
}

func (r *Remote) CombinedOutput(ctx context.Context) ([]byte, error) {
	return nil, nil
}
