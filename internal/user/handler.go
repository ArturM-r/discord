package user

import (
	"context"
	"discord/internal/errs"
	"discord/internal/jwt"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type Service interface {
	Registration(ctx context.Context, email string, password string) (*Response, error)
	Login(ctx context.Context, email string, password string) (*LoginResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var user Register

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jwt.JsonError(w, "bad request", http.StatusBadRequest)
		return
	}
	response, err := h.service.Registration(r.Context(), user.Email, user.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBadRequest):
			jwt.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errs.ErrEmailExists):
			jwt.JsonError(w, err.Error(), http.StatusConflict)
		default:
			jwt.JsonError(w, "registration failed", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user Login

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jwt.JsonError(w, "bad request", 400)
		return
	}
	response, err := h.service.Login(r.Context(), user.Email, user.Password)
	if err != nil {
		jwt.JsonError(w, err.Error(), 401)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}
