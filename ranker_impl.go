package main

import (
	"context"
	"fmt"

	"github.com/clinia/models-client-go/common"
	"github.com/clinia/models-client-go/datatype"
)

const (
	rankerQueryInputKey      string = "query"
	rankerQueryInputDatatype = datatype.Bytes

	// TODO: Change to text
	rankerPassageInputKey      string = "passage"
	rankerPassageInputDatatype = datatype.Bytes

	rankerScoreOutputKey string = "score"
)

type ranker struct {
	requester common.Requester
}

var _ Ranker = (*ranker)(nil)

func NewRanker(opts common.ClientOptions) Ranker {
	return &ranker{
		requester: opts.Requester,
	}
}

// Rank implements Ranker.
func (r *ranker) Rank(ctx context.Context, modelName string, modelVersion string, req RankRequest) (RankResponse, error) {
	// Duplicate query to be the same size as texts
	inputQueries := make([]string, len(req.Texts))
	for i := range req.Texts {
		inputQueries[i] = req.Query
	}

	// We don't specify the shape considering it calculated inside the infer function
	// when transforming the string content to the raw input.
	inputs := []common.Input{
		{
			Name:     rankerQueryInputKey,
			Datatype: rankerQueryInputDatatype,
			Contents: []common.Content{
				{
					StringContents: inputQueries,
				},
			},
		},
		{
			Name:     rankerPassageInputKey,
			Datatype: rankerPassageInputDatatype,
			Contents: []common.Content{
				{
					StringContents: req.Texts,
				},
			},
		},
	}

	// The ranker model has only one input and one output.
	outputKeys := []string{rankerScoreOutputKey}

	outputs, err := r.requester.Infer(ctx, modelName, modelVersion, inputs, outputKeys)

	if err != nil {
		return RankResponse{}, err
	}

	// Since we have only one output, we can directly access the first output.
	// We already check the size of the output in the infer function.
	scores := outputs[0].GetFp32Contents()

	// Flatten the 2D slice into a 1D slice
	var flattenedScores []float32
	for _, score := range scores {
		if len(score) != 1 {
			return RankResponse{}, fmt.Errorf("Expected a single score per passage, but got %d elements", len(score))
		}
		flattenedScores = append(flattenedScores, score...)
	}
	return RankResponse{
		ID:     req.ID,
		Scores: flattenedScores,
	}, nil
}
