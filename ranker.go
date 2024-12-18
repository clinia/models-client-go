package main

import "context"

type Ranker interface {
	// Rank returns the ranked results of the given texts.
	Rank(ctx context.Context, modelName, modelVersion string, req RankRequest) (RankResponse, error)
}

type RankRequest struct {
	ID       string
	Query    string
	Texts []string
}

type RankResponse struct {
	ID     string
	Scores []float32
}
