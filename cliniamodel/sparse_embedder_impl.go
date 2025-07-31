package cliniamodel

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

const (
	sparseEmbedderInputKey      string = "text"
	sparseEmbedderOutputKey     string = "embedding"
	sparseEmbedderInputDatatype        = datatype.Bytes
)

// embedder is a struct that implements the SparseEmbedder interface.
type sparseEmbedder struct {
	requester common.Requester
}

var _ SparseEmbedder = (*sparseEmbedder)(nil)

// NewEmbedder creates a new instance of embedder.
func NewSparseEmbedder(ctx context.Context, opts common.ClientOptions) SparseEmbedder {
	return &sparseEmbedder{
		requester: opts.Requester,
	}
}

// Embed generates embeddings for the given texts using the specified model and version.
func (e *sparseEmbedder) Embed(ctx context.Context, modelName, modelVersion string, req SparseEmbedRequest) (*SparseEmbedResponse, error) {
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

	flat := res.Outputs[0].Content.StringContents
	if len(flat) == 0 {
		return nil, errors.New("string matrix is empty")
	}

	embeddings := make([]map[string]float32, len(flat))
	for i, embJSON := range flat {
		var m map[string]float32
		if err := json.Unmarshal([]byte(embJSON), &m); err != nil {
			return nil, errors.New("unmarshal embedding " + strconv.Itoa(i) + ": " + err.Error())
		}
		embeddings[i] = m
	}

	return &SparseEmbedResponse{
		ID:         req.ID,
		Embeddings: embeddings,
	}, nil
}

// Ready implements the Embedder interface. It checks the readiness status of the model.
func (c *sparseEmbedder) Ready(ctx context.Context, modelName string, modelVersion string) error {
	return c.requester.Ready(ctx, modelName, modelVersion)
}
