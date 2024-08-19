package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	paymentService weaver.Ref[PaymentService]
	gatewayA       weaver.Ref[GatewayA]
	gatewayB       weaver.Ref[GatewayB]
	gatewayRouter  weaver.Ref[GatewayRouter]
	merger         weaver.Ref[Merger]
	lis            weaver.Listener
}

func run(ctx context.Context, app *app) error {
	fmt.Printf("Payment server listening on %v\n", app.lis)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/deposit", app.handleDeposit)
	r.Post("/withdrawal", app.handleWithdrawal)
	r.Post("/callback/{gateway}", app.handleCallback)
	r.Get("/transaction/{id}", app.handleGetTransaction)

	// serve the static file in /static
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./redoc-static.html")
	})

	return http.Serve(app.lis, r)
}

func (a *app) handleDeposit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount  float64 `json:"amount"`
		Gateway string  `json:"gateway"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err.Error(), "invalid_request", http.StatusBadRequest)
		return
	}

	tx, err := a.paymentService.Get().ProcessPayment(r.Context(), req.Amount, "deposit", req.Gateway)
	if err != nil {
		jsonError(w, err.Error(), "payment_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (a *app) handleWithdrawal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount  float64 `json:"amount"`
		Gateway string  `json:"gateway"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err.Error(), "invalid_request", http.StatusBadRequest)
		return
	}

	tx, err := a.paymentService.Get().ProcessPayment(r.Context(), req.Amount, "withdrawal", req.Gateway)
	if err != nil {
		jsonError(w, err.Error(), "payment_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (a *app) handleCallback(w http.ResponseWriter, r *http.Request) {
	gateway := chi.URLParam(r, "gateway")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		jsonError(w, err.Error(), "callback_error", http.StatusBadGateway)
		return
	}

	if err := a.paymentService.Get().HandleCallback(r.Context(), gateway, body); err != nil {
		jsonError(w, err.Error(), "callback_error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (a *app) handleGetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	tx, err := a.paymentService.Get().GetTransaction(r.Context(), id)
	if err != nil {
		jsonError(w, err.Error(), "transaction_not_found", http.StatusBadRequest)
		return
	}
	// this could also be used as an entirely different endpoint
	xmlJson, err := a.merger.Get().Merge(r.Context(), *tx)
	if err != nil {
		jsonError(w, err.Error(), "merger_error", http.StatusInternalServerError)
		return
	}
	log.Printf("the returned type is: %+v", xmlJson)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(xmlJson)
}
