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
	modelName := "chunker"
	modelVersion := "120252801110000"
	texts := []string{
		"Hello, how are you?",
		"The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.The AIRS Taxonomy (Taxonomy of Human Services) is a comprehensive and standardized classification system developed by the Alliance of Information and Referral Systems (AIRS). It is widely used to organize and categorize information about health, human, and social services. This taxonomy serves as a framework for Information & Referral (I&R) providers, enabling them to offer accurate and consistent referrals to services in their community.",
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

	chunker := cliniamodel.NewChunker(ctx, common.ClientOptions{
		Requester: requester,
	})

	req := cliniamodel.ChunkRequest{
		Texts: texts,
	}

	res, err := chunker.Chunk(ctx, modelName, modelVersion, req)
	if err != nil {
		log.Fatalf("chunk error: %v", err)
	}

	fmt.Println(res)
}
