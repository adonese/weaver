package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, amount float64, transactionType, gateway string) (*Transaction, error)
	HandleCallback(ctx context.Context, gateway string, data []byte) error
	GetTransaction(ctx context.Context, id string) (*Transaction, error)
}

type paymentService struct {
	weaver.Implements[PaymentService]
	store  sync.Map
	router weaver.Ref[GatewayRouter]
}

func (s *paymentService) ProcessPayment(ctx context.Context, amount float64, transactionType, gateway string) (*Transaction, error) {
	txID, err := s.router.Get().RoutePayment(ctx, amount, transactionType, gateway)
	if err != nil {
		return nil, err
	}

	tx := &Transaction{
		ID:        txID,
		Amount:    amount,
		Type:      transactionType,
		Status:    "pending",
		Gateway:   gateway,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.store.Store(txID, tx)
	return tx, nil
}

func (s *paymentService) HandleCallback(ctx context.Context, gateway string, data []byte) error {
	var txID, status string
	var err error

	switch gateway {
	case "gateway_a":
		var callback struct {
			TransactionID string `json:"transaction_id"`
			Status        string `json:"status"`
		}
		err = json.Unmarshal(data, &callback)
		txID, status = callback.TransactionID, callback.Status
	case "gateway_b":
		var callback struct {
			XMLName       xml.Name `xml:"callback"`
			TransactionID string   `xml:"transaction_id"`
			Status        string   `xml:"status"`
		}
		err = xml.Unmarshal(data, &callback)
		txID, status = callback.TransactionID, callback.Status
	default:
		// either fail this way, or try to fallback to json marshalling.
		return fmt.Errorf("unknown gateway: %s", gateway)
	}

	if err != nil {
		return err
	}

	// Update the transaction status
	if txInterface, ok := s.store.Load(txID); ok {
		tx := txInterface.(*Transaction)
		tx.Status = status
		tx.UpdatedAt = time.Now()
		s.store.Store(txID, tx)
	} else {
		return fmt.Errorf("transaction not found: %s", txID)
	}

	return nil
}

func (s *paymentService) GetTransaction(ctx context.Context, id string) (*Transaction, error) {
	tx, ok := s.store.Load(id)
	if !ok {
		return nil, errors.New("transaction not found")
	}
	return tx.(*Transaction), nil
}
