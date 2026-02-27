package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Xanaduxan/wallet-go/internal/service/auth"
)

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistrationResponse struct {
	AccessToken string `json:"access_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

var authService *auth.Service

func SetAuthService(s *auth.Service) { authService = s }

func Registration(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	token, err := authService.Registration(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrValidationFailed) {
			http.Error(w, "email already taken", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := RegistrationResponse{AccessToken: token}
	writeJSON(w, http.StatusCreated, resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := req.Validate(); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	token, err := authService.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	resp := LoginResponse{AccessToken: token}
	writeJSON(w, http.StatusOK, resp)
}

func DeleteMe(w http.ResponseWriter, r *http.Request) {

	emailValue := r.Context().Value("email")

	email, ok := emailValue.(string)
	if !ok || email == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := authService.DeleteByEmail(email)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateEmailPassword(email, password string) error {
	if email == "" || password == "" {
		return errors.New("email and password required")
	}
	return nil
}

func (req RegistrationRequest) Validate() error {
	return validateEmailPassword(req.Email, req.Password)
}

func (req LoginRequest) Validate() error {
	return validateEmailPassword(req.Email, req.Password)
}
