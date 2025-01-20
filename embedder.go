package main

import (
	"context"
)

type Embedder interface {
	// Embed returns the embeddings of the given texts.
	Embed(ctx context.Context, modelName, modelVersion string, req EmbedRequest) (EmbedResponse, error)
}

type EmbedRequest struct {
	ID    string
	Texts []string
}

type EmbedResponse struct {
	ID         string
	Embeddings [][]float32
}
