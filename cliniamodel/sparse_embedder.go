package cliniamodel

import "context"

type SparseEmbedder interface {
	// SparseEmbed returns the sparse embeddings of the given texts.
	SparseEmbed(ctx context.Context, modelName, modelVersion string, req SparseEmbedRequest) (*SparseEmbedResponse, error)
	// Ready checks if the model is ready to receive requests.
	Ready(ctx context.Context, modelName, modelVersion string) error
}

type SparseEmbedRequest struct {
	// ID is the unique identifier for the request.
	ID string
	// Texts is the list of texts to be embedded.
	Texts []string
}

type SparseEmbedResponse struct {
	// ID is the unique identifier for the response, corresponding to that of the request.
	ID string
	// Embeddings is the list of sparse embeddings for each text. Each embedding is a map
	// of tokens to their corresponding float values.
	Embeddings []map[string]float32
}
