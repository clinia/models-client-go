package main

import (
	"context"

	"github.com/clinia/models-client-go/common"
	"github.com/clinia/models-client-go/datatype"
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

func NewEmbedder(ctx context.Context, opts common.ClientOptions) (Embedder) {
	return &embedder{
		requester: opts.Requester,
	}
}

// Embed implements Embedder.
func (e *embedder) Embed(ctx context.Context, modelName, modelVersion string, req EmbedRequest) (EmbedResponse, error) {
	inputs := []common.Input{
		{
			Name:     embedderInputKey,
			Shape:    []int64{int64(len(req.Texts))},
			Datatype: embedderInputDatatype,
			Contents: []common.Content{
				{
					StringContents: req.Texts,
				},
			},
		},
	}

	// The embedder model has only one input and one output.
	outputKeys := []string{embedderOutputKey}

	outputs, err := e.requester.Infer(ctx, modelName, modelVersion, inputs, outputKeys)
	if err != nil {
		return EmbedResponse{}, err
	}

	// Since we have only one output, we can directly access the first output.
	// We already check the size of the output in the infer function therefore we can "safely" access the element 0.
	embeddings := outputs[0].GetFp32Contents()
	return EmbedResponse{
		ID:         req.ID,
		Embeddings: embeddings,
	}, nil
}
