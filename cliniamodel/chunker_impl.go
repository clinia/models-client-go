package cliniamodel

import (
	"context"
	"encoding/json"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/datatype"
)

const (
	chunkerInputKey      string = "text"
	chunkerOutputKey     string = "chunk"
	chunkerInputDatatype        = datatype.Bytes
	// Pad key is used to pad the ouput when the number of chunks across inputs is not the same.
	chunkerOutputPadKey = "pad"
)

// chunker is a struct that implements the Chunker interface.
type chunker struct {
	requester common.Requester
}

var _ Chunker = (*chunker)(nil)

// NewChunker creates a new chunker instance.
func NewChunker(ctx context.Context, opts common.ClientOptions) Chunker {
	return &chunker{
		requester: opts.Requester,
	}
}

// Chunk implements the Chunker interface. It takes a context, model name, model version, and a ChunkRequest as input,
// and returns a ChunkResponse or an error.
func (c *chunker) Chunk(ctx context.Context, modelName string, modelVersion string, req ChunkRequest) (*ChunkResponse, error) {
	// Prepare the inputs.
	inputs := []common.Input{
		{
			Name:     chunkerInputKey,
			Shape:    []int64{int64(len(req.Texts)), 1},
			Datatype: chunkerInputDatatype,
			Content: common.Content{
				StringContents: req.Texts,
			},
		},
	}

	outputKeys := []string{chunkerOutputKey}

	res, err := c.requester.Infer(ctx, common.InferRequest{
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
	// This is validated in the infer function and will be returned as an error if the output is not of the expected size.
	outputStringContents, err := res.Outputs[0].StringMatrixContent()
	chunks := [][]Chunk{}
	// We loop over the string contents and unmarshal them into the Chunk struct.
	for _, outputStringContent := range outputStringContents {
		textChunks := []Chunk{}
		for _, stringContent := range outputStringContent {
			// Skip the pad key.
			if stringContent == chunkerOutputPadKey {
				continue
			}

			var chunk Chunk
			if err := json.Unmarshal([]byte(stringContent), &chunk); err != nil {
				return nil, err
			}
			textChunks = append(textChunks, chunk)
		}
		chunks = append(chunks, textChunks)
	}

	return &ChunkResponse{
		ID:     res.ID,
		Chunks: chunks,
	}, nil
}
