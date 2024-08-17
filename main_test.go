package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ServiceWeaver/weaver/weavertest"
	"github.com/go-chi/chi/v5"
)

func TestAPIEndpoints(t *testing.T) {
	// Create a test runner
	runner := weavertest.Local
	runner.Test(t, func(t *testing.T, app *app) {
		// Create a new router and register the app's handlers
		r := chi.NewRouter()
		r.Post("/deposit", app.handleDeposit)
		r.Post("/withdrawal", app.handleWithdrawal)
		r.Post("/callback/{gateway}", app.handleCallback)
		r.Get("/transaction/{id}", app.handleGetTransaction)

		var depositTxID string
		var withdrawalID string

		// Test deposit
		t.Run("Deposit", func(t *testing.T) {
			reqBody := []byte(`{"amount": 100.0, "gateway": "gateway_a"}`)
			req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBuffer(reqBody))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}

			var tx Transaction
			err := json.Unmarshal(res.Body.Bytes(), &tx)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tx.Type != "deposit" || tx.Amount != 100.0 || tx.Gateway != "gateway_a" {
				t.Errorf("Unexpected transaction details: %+v", tx)
			}

			depositTxID = tx.ID
		})

		// Test withdrawal
		t.Run("Withdrawal", func(t *testing.T) {
			reqBody := []byte(`{"amount": 50.0, "gateway": "gateway_b"}`)
			req := httptest.NewRequest(http.MethodPost, "/withdrawal", bytes.NewBuffer(reqBody))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}

			var tx Transaction
			err := json.Unmarshal(res.Body.Bytes(), &tx)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tx.Type != "withdrawal" || tx.Amount != 50.0 || tx.Gateway != "gateway_b" {
				t.Errorf("Unexpected transaction details: %+v", tx)
			}
			withdrawalID = tx.ID
		})

		// Test callback
		t.Run("Callback", func(t *testing.T) {
			reqBody := []byte(`{"transaction_id": "` + depositTxID + `", "status": "completed"}`)
			req := httptest.NewRequest(http.MethodPost, "/callback/gateway_a", bytes.NewBuffer(reqBody))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}
		})

		// Test callback with XML (for gatewayB)
		t.Run("Callback XML", func(t *testing.T) {
			xmlBody := fmt.Sprintf(`
				<callback>
					<transaction_id>%s</transaction_id>
					<status>completed</status>
				</callback>
			`, withdrawalID)
			req := httptest.NewRequest(http.MethodPost, "/callback/gateway_b", bytes.NewBufferString(xmlBody))
			req.Header.Set("Content-Type", "application/xml")
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}
		})

		// Test *failed* withdrawal callback
		t.Run("Callback", func(t *testing.T) {
			reqBody := []byte(`{"transaction_id": "` + withdrawalID + `", "status": "failed"}`)
			req := httptest.NewRequest(http.MethodPost, "/callback/gateway_a", bytes.NewBuffer(reqBody))
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}
		})
		// Test get transaction
		t.Run("GetTransaction", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transaction/%s", depositTxID), nil)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}

			var mergedTx Response
			err := json.Unmarshal(res.Body.Bytes(), &mergedTx)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
			log.Printf("the transaction in GetTransaction is: %+v", mergedTx)

			if mergedTx.ID != depositTxID || mergedTx.Amount != 100.0 || mergedTx.Type != "deposit" || mergedTx.Gateway != "gateway_a" || mergedTx.Status != "completed" {
				t.Errorf("Unexpected transaction details: %+v", mergedTx)
			}

			if mergedTx.XMLResponse != "i am xml response" {
				t.Errorf("Expected extra_info in merged data")
			}
		})

		// GetTransaction for the failed withdrawal
		t.Run("GetTransaction-FailedWithdrawal", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transaction/%s", withdrawalID), nil)
			res := httptest.NewRecorder()

			r.ServeHTTP(res, req)

			if res.Code != http.StatusOK {
				t.Errorf("Expected status OK; got %v", res.Code)
			}

			var mergedTx Response
			err := json.Unmarshal(res.Body.Bytes(), &mergedTx)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}
			log.Printf("the transaction in GetTransaction is: %+v", mergedTx)

			if mergedTx.ID != withdrawalID || mergedTx.Amount != 50.0 || mergedTx.Type != "withdrawal" || mergedTx.Gateway != "gateway_b" || mergedTx.Status != "failed" {
				t.Errorf("Unexpected transaction details: %+v", mergedTx)
			}

			if mergedTx.XMLResponse != "i am xml response" {
				t.Errorf("Expected extra_info in merged data")
			}
		})
	})
}
