package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ServiceWeaver/weaver"
)

// Response our generic response with weaver marshalling. This one in particular is used to merge transactions.
type Response struct {
	weaver.AutoMarshal
	Transaction
	XMLResponse string
}

// ErrorResponse for all of our endpoints. We follow message - code structure, message describes the error occurred, while code can be used as i18n key to display the error in the user's language.
type ErrorResponse struct {
	weaver.AutoMarshal
	Message string
	Code    string
}

type Transaction struct {
	weaver.AutoMarshal
	ID     string
	Amount float64
	Type   string // "deposit" or "withdrawal"

	Status    string
	Gateway   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// jsonError is a helper function to return JSON-formatted error responses
func jsonError(w http.ResponseWriter, message string, code string, resCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: message, Code: code})
}
