package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
)

type Merger interface {
	Merge(ctx context.Context, jsonResonse Transaction) (Response, error)
}

type merger struct {
	weaver.Implements[Merger]
}

func (m *merger) Merge(_ context.Context, jsonResonse Transaction) (Response, error) {
	// TODO: Implement the logic to merge the transactions and return the response
	response := Response{
		Transaction: jsonResonse,
		XMLResponse: "i am xml response",
	}
	return response, nil
}

var _ weaver.NotRetriable = Merger.Merge
