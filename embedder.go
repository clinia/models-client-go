package main

import "context"

type Embedder interface {
	Embed(ctx context.Context) error
}
