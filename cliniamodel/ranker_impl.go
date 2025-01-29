package cliniamodel

import (
	"context"
	"fmt"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

const (
	rankerQueryInputKey      string            = "query"
	rankerQueryInputDatatype datatype.Datatype = datatype.Bytes

	// TODO: Change to text
	rankerPassageInputKey      string            = "passage"
	rankerPassageInputDatatype datatype.Datatype = datatype.Bytes

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
func (r *ranker) Rank(ctx context.Context, modelName string, modelVersion string, req RankRequest) (*RankResponse, error) {
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
			Content: common.Content{
				StringContents: inputQueries,
			},
		},
		{
			Name:     rankerPassageInputKey,
			Datatype: rankerPassageInputDatatype,
			Content: common.Content{
				StringContents: req.Texts,
			},
		},
	}

	// The ranker model has only one input and one output.
	outputKeys := []string{rankerScoreOutputKey}

	res, err := r.requester.Infer(ctx, common.InferRequest{
		ID:           req.ID,
		ModelName:    modelName,
		ModelVersion: modelVersion,
		Inputs:       inputs,
		OutputKeys:   outputKeys,
	})

	if err != nil {
		return nil, err
	}

	// Since we have only one output, we can directly access the first output.
	// We already check the size of the output in the infer function.
	scores, err := res.Outputs[0].Fp32MatrixContent()
	if err != nil {
		return nil, err
	}

	// Flatten the 2D slice into a 1D slice
	var flattenedScores []float32
	for _, score := range scores {
		if len(score) != 1 {
			return nil, fmt.Errorf("Expected a single score per passage, but got %d elements", len(score))
		}
		flattenedScores = append(flattenedScores, score...)
	}
	return &RankResponse{
		ID:     res.ID,
		Scores: flattenedScores,
	}, nil
}
