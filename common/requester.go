package common

import (
	"context"
)

type Requester interface {
	// 
	Infer(ctx context.Context, modelName, modelVersion string, inputs []Input, outputKeys []string) ([]Output, error)
	Stream(ctx context.Context, modelName, modelVersion string, inputs []Input) (chan <- string,  error)
	Close() error
}


type RequesterConfig struct {
	Host Host
}


