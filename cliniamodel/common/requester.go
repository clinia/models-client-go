package common

import (
	"context"
)

type Requester interface {
	// Infer sends a request to the model server to perform inference on the given inputs.
	Infer(ctx context.Context, req InferRequest) (*InferResponse, error)
	// Stream sends a request to the model server to perform inference on the given inputs and returns a channel to stream the output.
	Stream(ctx context.Context, modelName, modelVersion string, inputs []Input) (chan<- string, error)
	// Close closes the connection to the model server.
	Close() error
}

type RequesterConfig struct {
	Host Host
}

type InferRequest struct {
	// ID is a unique identifier for the request.
	ID string
	// ModelName is the name of the model to be used for inference.
	ModelName string
	// ModelVersion is the version of the model to be used for inference.
	ModelVersion string
	// Inputs will be a list of inputs to be used for inference.
	Inputs []Input
	// OutputKeys will be a list of keys to be used to access the outputs.
	OutputKeys []string
}

type InferResponse struct {
	// ID is a unique identifier for the response.
	ID string
	// Outputs will be a list of outputs for the given inputs.
	Outputs []Output
}
