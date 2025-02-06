package main

import (
	"context"
	"fmt"
	"log"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"

	"github.com/clinia/models-client-go/cliniamodel"
)

func main() {
	modelName := "ranker_medical_journals_qa"
	modelVersion := "120240905185925"
	query := "hello, how are you?"
	texts := []string{
		"Clinia is based in Montreal",
	}

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

	ranker := cliniamodel.NewRanker(common.ClientOptions{
		Requester: requester,
	})

	req := cliniamodel.RankRequest{
		Texts: texts,
		Query: query,
	}

	res, err := ranker.Rank(ctx, modelName, modelVersion, req)
	if err != nil {
		log.Fatalf("rank error: %v", err)
	}

	fmt.Println(res)
}
