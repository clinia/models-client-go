package cliniamodel

import (
	"context"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

const (
	embedderInputKey      string = "text"
	embedderOutputKey     string = "embedding"
	embedderInputDatatype        = datatype.Bytes
)

type embedder struct {
	requester common.Requester
}

var _ Embedder = (*embedder)(nil)

func NewEmbedder(ctx context.Context, opts common.ClientOptions) Embedder {
	return &embedder{
		requester: opts.Requester,
	}
}

// Embed implements Embedder.
func (e *embedder) Embed(ctx context.Context, modelName, modelVersion string, req EmbedRequest) (*EmbedResponse, error) {
	inputs := []common.Input{
		{
			Name:     embedderInputKey,
			Shape:    []int64{int64(len(req.Texts))},
			Datatype: embedderInputDatatype,
			Content: common.Content{
				StringContents: req.Texts,
			},
		},
	}

	// The embedder model has only one input and one output.
	outputKeys := []string{embedderOutputKey}

	res, err := e.requester.Infer(ctx, common.InferRequest{
		ID:           req.ID,
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Inputs:       inputs,
		OutputKeys:   outputKeys,
	})
	if err != nil {
		return nil, err
	}

	// Since we have only one output, we can directly access the first output.
	// We already check the size of the output in the infer function therefore we can "safely" access the element 0.
	embeddings, err := res.Outputs[0].Fp32MatrixContent()
	if err != nil {
		return nil, err
	}

	return &EmbedResponse{
		ID:         res.ID,
		Embeddings: embeddings,
	}, nil
}
