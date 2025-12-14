package protocols

import "context"

type PipelineStep interface {
	Name() string
	Process(ctx context.Context, data []byte) ([]byte, error)
}
