package requestergrpc

import (
	"context"
	"fmt"

	"github.com/clinia/models-client-go/common"
	"github.com/clinia/models-client-go/datatype"
	"github.com/clinia/models-client-go/requestergrpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type requester struct{
	conn *grpc.ClientConn

	inferenceServiceClient gen.GRPCInferenceServiceClient
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
		inferenceServiceClient: gen.NewGRPCInferenceServiceClient(conn),
	}, nil
}

// Infer implements common.Requester.
func (r *requester) Infer(ctx context.Context, modelName string, modelVersion string, inputs []common.Input, outputKeys []string) ([]common.Output, error) {
	// Prepare input tensors
	grpcInputs := make([]*gen.ModelInferRequest_InferInputTensor, len(inputs))
	rawInputs := make([][]byte, len(inputs))
	for i, input := range inputs {
		
		// For now, we only support bytes/string datatype
		// TODO: Support other datatypes
		if input.Datatype != datatype.Bytes {
			return nil, fmt.Errorf("unsupported datatype: %v", input.Datatype)
		}

		rawInputContents, shape, err := preprocess(input.GetStringContents()[i])
		if err != nil {
			return nil, err
		}

		grpcInputs[i] = &gen.ModelInferRequest_InferInputTensor{
			Name: input.Name,
			Shape: shape,
			Datatype: string(input.Datatype),
		}

		rawInputs[i] = rawInputContents
	}

	// Prepare output keys
	grpcOutputs := make([]*gen.ModelInferRequest_InferRequestedOutputTensor, len(inputs))
	for i, outputKey := range outputKeys {
		grpcOutputs[i] = &gen.ModelInferRequest_InferRequestedOutputTensor{
			Name: outputKey,
		}
	}

	res, err := r.inferenceServiceClient.ModelInfer(ctx, &gen.ModelInferRequest{
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Inputs: 		 grpcInputs,
		Outputs: grpcOutputs,
		RawInputContents: rawInputs,
	})

	if err != nil {
		return nil, err
	}

	// TODO: Check resp ID
	if res.Id != "" {
		return nil, fmt.Errorf("unexpected response ID: %s", res.Id)
	}

	// Check if the number of output keys matches the number of outputs
	if len(res.RawOutputContents) != len(outputKeys) {
		return nil, fmt.Errorf("expected %d output keys, got %d", len(outputKeys), len(res.RawOutputContents))
	}

	// Prepare output tensors
	outputs := make([]common.Output, len(res.Outputs))
	for i, rawOutput := range res.RawOutputContents {
		resOutput := res.Outputs[i]

		// For now, we only support FP32 datatype
		// TODO: Support other datatypes
		if resOutput.Datatype != string(datatype.Fp32) {
			return nil, fmt.Errorf("unsupported output datatype: %v", resOutput.Datatype)
		}

		output32, err := postprocessFp32(rawOutput, resOutput.Shape)
		if err != nil {
			return nil, err
		}

		contents := make([]common.Content, len(output32))
		for i, v := range output32 {
			contents[i] = common.Content{Fp32Contents: v}
		}

		outputs[i] = common.Output{
			Name: resOutput.Name,
			Shape: resOutput.Shape,
			Datatype: datatype.Fp32,
			Contents: contents,
		}
	}

	return outputs, nil
}

// Stream implements common.Requester.
func (r *requester) Stream(ctx context.Context, modelName string, modelVersion string, inputs []common.Input) (chan<- string, error) {
	panic("unimplemented")
}

func (r *requester) Close() error {
	return r.conn.Close()
}
