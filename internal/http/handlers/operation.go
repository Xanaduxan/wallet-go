package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Xanaduxan/wallet-go/internal/service/operations"
	"github.com/google/uuid"
)

type CreateOperationRequest struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

type CreateOperationResponse struct {
	ID uuid.UUID `json:"id"`
}

var operationService *operations.Service

func SetOperationService(s *operations.Service) { operationService = s }

func CreateOperation(w http.ResponseWriter, r *http.Request) {
	var req CreateOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	email, _ := r.Context().Value("email").(string)
	if email == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := operationService.CreateOperation(email, req.Name, req.Type, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, operations.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
			return
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusCreated, CreateOperationResponse{ID: id})
}

func (req CreateOperationRequest) Validate() error {
	if req.Name == "" ||
		(req.Type != "income" && req.Type != "expense") ||
		req.Amount <= 0 {
		return operations.ErrInvalidInput
	}
	return nil
}
