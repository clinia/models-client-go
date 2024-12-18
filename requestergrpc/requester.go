package requestergrpc

import (
	"context"

	"github.com/clinia/models-client-go/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type requester struct{
	conn *grpc.ClientConn

	inferenceServiceClient GRPCInferenceServiceClient
}

var _ common.Requester = (*requester)(nil)

func NewRequester(ctx context.Context, cfg common.RequesterConfig) (common.Requester, error) {
	opts := []grpc.DialOption{}

	// Set insecure credentials if the host is HTTP
	if cfg.Host.Scheme == common.HTTP {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(cfg.Host.Host(), opts...)
	if err != nil {
		return nil, err
	}

	return &requester{
		conn: conn,
		inferenceServiceClient: NewGRPCInferenceServiceClient(conn),
	}, nil
}

// Infer implements common.Requester.
func (r *requester) Infer(ctx context.Context, modelName string, modelVersion string, inputs []common.Input) ([]common.Output, error) {
	panic("unimplemented")
}

// Stream implements common.Requester.
func (r *requester) Stream(ctx context.Context, modelName string, modelVersion string, inputs []common.Input) (chan<- string, error) {
	panic("unimplemented")
}

func (r *requester) Close() error {
	return r.conn.Close()
}
