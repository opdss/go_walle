package command

import "context"

type Command interface {
	Run(ctx context.Context) error
	CombinedOutput(ctx context.Context) ([]byte, error)
}
