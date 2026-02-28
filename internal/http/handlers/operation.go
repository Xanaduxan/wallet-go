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

type UpdateOperationRequest struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount"`
}

var operationService *operations.Service

func SetOperationService(s *operations.Service) { operationService = s }

func GetOperation(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value("email").(string)
	if email == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	opID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	op, err := operationService.GetOperation(email, opID)
	if err != nil {
		switch {
		case errors.Is(err, operations.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		case errors.Is(err, operations.ErrForbidden):
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusOK, op)
}

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
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, http.StatusCreated, CreateOperationResponse{ID: id})
}

func DeleteOperation(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value("email").(string)
	if email == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	opID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := operationService.DeleteOperation(email, opID); err != nil {
		switch {
		case errors.Is(err, operations.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
		case errors.Is(err, operations.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		case errors.Is(err, operations.ErrForbidden):
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateOperation(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value("email").(string)
	if email == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	opID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req UpdateOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if err := operationService.UpdateOperation(email, opID, req.Name, req.Type, req.Amount); err != nil {
		switch {
		case errors.Is(err, operations.ErrInvalidInput):
			http.Error(w, "invalid input", http.StatusBadRequest)
		case errors.Is(err, operations.ErrNotFound):
			http.Error(w, "not found", http.StatusNotFound)
		case errors.Is(err, operations.ErrForbidden):
			http.Error(w, "forbidden", http.StatusForbidden)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (req CreateOperationRequest) Validate() error {
	return validateOperation(req.Name, req.Type, req.Amount)
}

func (req UpdateOperationRequest) Validate() error {
	return validateOperation(req.Name, req.Type, req.Amount)
}
