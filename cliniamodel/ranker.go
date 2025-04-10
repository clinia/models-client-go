package cliniamodel

import "context"

type Ranker interface {
	// Rank returns the ranked results of the given texts.
	Rank(ctx context.Context, modelName, modelVersion string, req RankRequest) (*RankResponse, error)
	// Ready checks if the model is ready to receive requests.
	Ready(ctx context.Context, modelName, modelVersion string) error
}

type RankRequest struct {
	// ID is the unique identifier for the request. If not provided, a random UUIDv4 will be generated.
	ID string
	// Query is the query to rank the passages against.
	Query string
	// Texts is the list of passages to be ranked.
	Texts []string
}

type RankResponse struct {
	// ID is the unique identifier for the response, corresponding to that of the request.
	ID string
	// Scores is the list of scores for each pair of query and passage.
	Scores []float32
}
