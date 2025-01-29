package cliniamodel

import "context"

type Chunker interface {
	Chunk(ctx context.Context, modelName, modelVersion string, req ChunkRequest) (*ChunkResponse, error)
}

type ChunkRequest struct {
	// ID is a unique identifier for the request.
	ID string
	// Texts will be a list of texts to be chunked.
	Texts []string
}

type ChunkResponse struct {
	// ID is a unique identifier for the response.
	ID string
	// Chunks will be a list of chunks for each text in the input.
	// The length of the list will be equal to the length of the input texts.
	Chunks [][]Chunk
}

type Chunk struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	StartIndex int    `json:"startIndex"`
	EndIndex   int    `json:"endIndex"`
	TokenCount int    `json:"tokenCount"`
}
