package requestergrpc

import (
	"context"
	"fmt"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
	requestergrpc "github.com/clinia/models-client-go/cliniamodel/requestergrpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type requester struct {
	conn *grpc.ClientConn

	inferenceServiceClient requestergrpc.GRPCInferenceServiceClient
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
		conn:                   conn,
		inferenceServiceClient: requestergrpc.NewGRPCInferenceServiceClient(conn),
	}, nil
}

// Infer implements common.Requester.
func (r *requester) Infer(ctx context.Context, req common.InferRequest) (*common.InferResponse, error) {
	// Prepare input tensors
	grpcInputs := make([]*requestergrpc.ModelInferRequest_InferInputTensor, len(req.Inputs))
	rawInputs := make([][]byte, len(req.Inputs))
	for i, input := range req.Inputs {

		// For now, we only support bytes/string datatype
		if input.Datatype != datatype.Bytes {
			return nil, fmt.Errorf("unsupported datatype: %v", input.Datatype)
		}

		rawInputContents, shape, err := preprocessString(input.GetStringContents())
		if err != nil {
			return nil, err
		}

		grpcInputs[i] = &requestergrpc.ModelInferRequest_InferInputTensor{
			Name:     input.Name,
			Shape:    shape,
			Datatype: string(input.Datatype),
		}

		rawInputs[i] = rawInputContents
	}

	// Prepare output keys
	grpcOutputs := make([]*requestergrpc.ModelInferRequest_InferRequestedOutputTensor, len(req.OutputKeys))
	for i, outputKey := range req.OutputKeys {
		grpcOutputs[i] = &requestergrpc.ModelInferRequest_InferRequestedOutputTensor{
			Name: outputKey,
		}
	}

	res, err := r.inferenceServiceClient.ModelInfer(ctx, &requestergrpc.ModelInferRequest{
		Id:               req.ID,
		ModelName:        req.ModelName,
		ModelVersion:     req.ModelVersion,
		Inputs:           grpcInputs,
		Outputs:          grpcOutputs,
		RawInputContents: rawInputs,
	})

	if err != nil {
		return nil, err
	}

	if res.Id != req.ID {
		return nil, fmt.Errorf("unexpected response ID: %s", res.Id)
	}

	// Check if the number of output keys matches the number of outputs
	if len(res.RawOutputContents) != len(req.OutputKeys) {
		return nil, fmt.Errorf("expected %d output keys, got %d", len(req.OutputKeys), len(res.RawOutputContents))
	}

	// Prepare output tensors
	outputs := make([]common.Output, len(res.Outputs))
	for i, rawOutput := range res.RawOutputContents {
		resOutput := res.Outputs[i]

		// For now, we only support FP32 & bytes/string datatype
		switch resOutput.Datatype {
		case string(datatype.Fp32):
			fp32Contents, err := decodeFloat32(rawOutput)
			if err != nil {
				return nil, err
			}

			outputs[i] = common.Output{
				Name:     resOutput.Name,
				Shape:    resOutput.Shape,
				Datatype: datatype.Fp32,
				Content: common.Content{
					Fp32Contents: fp32Contents,
				},
			}
		case string(datatype.Bytes):
			stringContents, err := decodeString(rawOutput)
			if err != nil {
				return nil, err
			}

			outputs[i] = common.Output{
				Name:     resOutput.Name,
				Shape:    resOutput.Shape,
				Datatype: datatype.Bytes,
				Content: common.Content{
					StringContents: stringContents,
				},
			}
		default:
			return nil, fmt.Errorf("unsupported output datatype: %v", resOutput.Datatype)
		}
	}

	return &common.InferResponse{
		ID:      res.Id,
		Outputs: outputs,
	}, nil
}

// Stream implements common.Requester.
func (r *requester) Stream(ctx context.Context, modelName string, modelVersion string, inputs []common.Input) (chan<- string, error) {
	panic("unimplemented")
}

func (r *requester) Close() error {
	return r.conn.Close()
}
