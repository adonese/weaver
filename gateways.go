package main

import (
	"context"

	"github.com/ServiceWeaver/weaver"
	"github.com/google/uuid"
)

type GatewayA interface {
	ProcessDeposit(ctx context.Context, amount float64) (string, error)
	ProcessWithdrawal(ctx context.Context, amount float64) (string, error)
}

type gatewayA struct {
	weaver.Implements[GatewayA]
}

func (g *gatewayA) ProcessDeposit(ctx context.Context, amount float64) (string, error) {
	return uuid.New().String(), nil
}

func (g *gatewayA) ProcessWithdrawal(ctx context.Context, amount float64) (string, error) {
	return uuid.New().String(), nil
}

type GatewayB interface {
	ProcessDeposit(ctx context.Context, amount float64) (string, error)
	ProcessWithdrawal(ctx context.Context, amount float64) (string, error)
}

type gatewayB struct {
	weaver.Implements[GatewayB]
}

func (g *gatewayB) ProcessDeposit(ctx context.Context, amount float64) (string, error) {
	return uuid.New().String(), nil
}

func (g *gatewayB) ProcessWithdrawal(ctx context.Context, amount float64) (string, error) {
	return uuid.New().String(), nil
}
