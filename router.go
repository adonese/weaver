package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
)

type GatewayRouter interface {
	RoutePayment(ctx context.Context, amount float64, transactionType, gateway string) (string, error)
}

type gatewayRouter struct {
	weaver.Implements[GatewayRouter]
	gatewayA weaver.Ref[GatewayA]
	gatewayB weaver.Ref[GatewayB]
}

// RoutePayment manages routing traffic between
func (r *gatewayRouter) RoutePayment(ctx context.Context, amount float64, transactionType, gateway string) (string, error) {
	var txID string
	var err error

	switch gateway {
	case "gateway_a":
		g := r.gatewayA.Get()
		if transactionType == "deposit" {
			txID, err = g.ProcessDeposit(ctx, amount)
		} else {
			txID, err = g.ProcessWithdrawal(ctx, amount)
		}
	case "gateway_b":
		g := r.gatewayB.Get()
		if transactionType == "deposit" {
			txID, err = g.ProcessDeposit(ctx, amount)
		} else {
			txID, err = g.ProcessWithdrawal(ctx, amount)
		}
	default:
		return "", fmt.Errorf("unknown gateway: %s", gateway)
	}

	return txID, err
}
