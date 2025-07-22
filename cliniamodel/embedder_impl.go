package cliniamodel

import (
	"context"
	"errors"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

const (
	embedderInputKey      string = "text"
	embedderOutputKey     string = "embedding"
	embedderInputDatatype        = datatype.Bytes
)

// embedder is a struct that implements the Embedder interface.
type embedder struct {
	requester common.Requester
}

var _ Embedder = (*embedder)(nil)

// NewEmbedder creates a new instance of embedder.
func NewEmbedder(ctx context.Context, opts common.ClientOptions) Embedder {
	return &embedder{
		requester: opts.Requester,
	}
}

// Embed generates embeddings for the given texts using the specified model and version.
func (e *embedder) Embed(ctx context.Context, modelName, modelVersion string, req EmbedRequest) (*EmbedResponse, error) {
	if len(req.Texts) == 0 {
		return nil, errors.New("texts cannot be empty")
	}

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

// Ready implements the Embedder interface. It checks the readiness status of the model.
func (c *embedder) Ready(ctx context.Context, modelName string, modelVersion string) error {
	return c.requester.Ready(ctx, modelName, modelVersion)
}
