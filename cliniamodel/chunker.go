package cliniamodel

import "context"

type Chunker interface {
	// Chunk returns the chunked results of the given texts.
	Chunk(ctx context.Context, modelName, modelVersion string, req ChunkRequest) (*ChunkResponse, error)
	// Ready checks if the model is ready to receive requests.
	Ready(ctx context.Context, modelName, modelVersion string) error
}

type ChunkRequest struct {
	// ID is the unique identifier for the request.
	ID string
	// Texts is the list of texts to be chunked.
	Texts []string
}

type ChunkResponse struct {
	// ID is the unique identifier for the response, corresponding to that of the request.
	ID string
	// Chunks is the list of chunks in which each text is split. The outer list length matches the number of input texts.
	Chunks [][]Chunk
}

type Chunk struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	TokenCount int    `json:"tokenCount"`
}
