package cliniamodel

import (
	"context"
)

type Embedder interface {
	// Embed returns the embeddings of the given texts.
	Embed(ctx context.Context, modelName, modelVersion string, req EmbedRequest) (*EmbedResponse, error)
}

type EmbedRequest struct {
	// ID is the unique identifier for the request.
	ID string
	// Texts is the list of texts to be embedded.
	Texts []string
}

type EmbedResponse struct {
	// ID is the unique identifier for the response, corresponding to that of the request.
	ID string
	// Embeddings is the list of embeddings for each text. Each embedding is a list of floats, corresponding to the embedding dimensions. The outer list length matches the number of input texts.
	Embeddings [][]float32
}
