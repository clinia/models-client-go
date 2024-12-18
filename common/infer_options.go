package common

type InferOptions struct {
	OutputKeys []string
}

type InferOption func(*InferOptions)

func WithOutputKeys(outputKeys []string) func(*InferOptions) {
	return func(o *InferOptions) {
		o.OutputKeys = outputKeys
	}
}
