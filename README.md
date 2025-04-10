# Clinia Models Client Go

## Features

- Native Go implementation for communicating with Clinia's models on NVIDIA Triton Inference Server
- gRPC-based communication for efficient and reliable model inference
- Support for concurrent requests
- Type-safe interfaces for all models and requests
- Support for batch processing for optimal performance
- Minimal dependencies and lightweight implementation
- Designed for internal Clinia services with VPC access

## Getting Started

To install the package, use go get:

```bash
go get github.com/clinia/models-client-go
```

You can now import the Clinia Models client in your project and use it to interact with the models.

## Playground Examples

### Embedder Model

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
)

func main() {
	// Get model name and version from environment variables.
	modelName := os.Getenv("CLINIA_MODEL_NAME")
	modelVersion := os.Getenv("CLINIA_MODEL_VERSION")
	if modelName == "" || modelVersion == "" {
		log.Fatal("Environment variables CLINIA_MODEL_NAME and CLINIA_MODEL_VERSION must be set")
	}

	// Define the texts to embed
	texts := []string{
		"What is the treatment for diabetes?",
		"What are the symptoms of heart disease?",
	}

	// Create a new requester with the host configuration.
	ctx := context.Background()
	requester, err := requestergrpc.NewRequester(ctx, common.RequesterConfig{
		Host: common.Host{
			Url:    "127.0.0.1",
			Port:   8001,
			Scheme: common.HTTP,
		},
	})
	if err != nil {
		log.Fatalf("failed to create requester: %v", err)
	}
	defer requester.Close()

	// Check if the server is ready
	if err := requester.Health(ctx); err != nil {
		log.Fatalf("server readiness check failed: %v", err)
	}

	// Create a new Embedder using the requester.
	embedder := cliniamodel.NewEmbedder(common.ClientOptions{
		Requester: requester,
	})

	// Check if the model is ready
	if err := embedder.Ready(ctx, modelName, modelVersion); err != nil {
		log.Fatalf("model readiness check failed: %v", err)
	}

	// Create an EmbedRequest with a generated ID and the texts.
	req := cliniamodel.EmbedRequest{
		ID:    uuid.New().String(),
		Texts: texts,
	}

	// Execute the embedding request.
	res, err := embedder.Embed(ctx, modelName, modelVersion, req)
	if err != nil {
		log.Fatalf("embed error: %v", err)
	}

	// Print results
	for i, embedding := range res.Embeddings {
		fmt.Printf("Text: %s\nEmbedding dimensions: %d\nFirst 5 values: %v\n\n",
			texts[i],
			len(embedding),
			embedding[:5],
		)
	}
}
```

### Ranker Model

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
)

func main() {
	// Get model name and version from environment variables.
	modelName := os.Getenv("CLINIA_MODEL_NAME")
	modelVersion := os.Getenv("CLINIA_MODEL_VERSION")
	if modelName == "" || modelVersion == "" {
		log.Fatal("Environment variables CLINIA_MODEL_NAME and CLINIA_MODEL_VERSION must be set")
	}

	// Define the query and texts
	query := "Where is Clinia based?"
	texts := []string{"Clinia is based in Montreal"}

	// Create a new requester with the host configuration.
	ctx := context.Background()
	requester, err := requestergrpc.NewRequester(ctx, common.RequesterConfig{
		Host: common.Host{
			Url:    "127.0.0.1",
			Port:   8001,
			Scheme: common.HTTP,
		},
	})
	if err != nil {
		log.Fatalf("failed to create requester: %v", err)
	}
	defer requester.Close()

	// Check if the server is ready
	if err := requester.Health(ctx); err != nil {
		log.Fatalf("server readiness check failed: %v", err)
	}

	// Create a new Ranker using the requester.
	ranker := cliniamodel.NewRanker(common.ClientOptions{
		Requester: requester,
	})

	// Check if the model is ready
	if err := ranker.Ready(ctx, modelName, modelVersion); err != nil {
		log.Fatalf("model readiness check failed: %v", err)
	}

	// Create a RankRequest with a generated ID, the query, and texts.
	req := cliniamodel.RankRequest{
		ID:    uuid.New().String(),
		Query: query,
		Texts: texts,
	}

	// Execute the ranking request.
	res, err := ranker.Rank(ctx, modelName, modelVersion, req)
	if err != nil {
		log.Fatalf("rank error: %v", err)
	}

	// Print results
	fmt.Println("rank result:", res)
}
```

### Chunker Model

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/clinia/models-client-go/cliniamodel"
	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
)

func main() {
	// Get model name and version from environment variables.
	modelName := os.Getenv("CLINIA_MODEL_NAME")
	modelVersion := os.Getenv("CLINIA_MODEL_VERSION")
	if modelName == "" || modelVersion == "" {
		log.Fatal("Environment variables CLINIA_MODEL_NAME and CLINIA_MODEL_VERSION must be set")
	}

	// Define the texts
	texts := []string{
		"This is a short text",
		"This is a longer text that. contains. mutliple sentences. but. should still be a single chunk",
	}

	// Create a new requester with the host configuration.
	ctx := context.Background()
	requester, err := requestergrpc.NewRequester(ctx, common.RequesterConfig{
		Host: common.Host{
			Url:    "127.0.0.1",
			Port:   8001,
			Scheme: common.HTTP,
		},
	})
	if err != nil {
		log.Fatalf("failed to create requester: %v", err)
	}
	defer requester.Close()

	// Check if the server is ready
	if err := requester.Health(ctx); err != nil {
		log.Fatalf("server readiness check failed: %v", err)
	}

	// Create a new Chunker using the requester.
	chunker := cliniamodel.NewChunker(common.ClientOptions{
		Requester: requester,
	})

	// Check if the model is ready
	if err := chunker.Ready(ctx, modelName, modelVersion); err != nil {
		log.Fatalf("model readiness check failed: %v", err)
	}

	// Create a ChunkRequest with a generated ID and the texts.
	req := cliniamodel.ChunkRequest{
		ID:    uuid.New().String(),
		Texts: texts,
	}

	// Execute the chunking request.
	res, err := chunker.Chunk(ctx, modelName, modelVersion, req)
	if err != nil {
		log.Fatalf("chunk error: %v", err)
	}

	// Print results
	fmt.Println("chunk result:", res)
}
```

## Note

This repository is automatically generated from a private repository within Clinia that contains additional resources including tests, mock servers, and development tools.

The version numbers of this package correspond to the same versions in the respective Python, Go and TypeScript public repositories, ensuring consistency across all implementations.

## License

Clinia Models Client Go is an open-sourced software licensed under the [MIT license](LICENSE).
