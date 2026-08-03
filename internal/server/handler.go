package server

import (
	"context"
	"discord/internal/jwt"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Service interface {
	CreateSrvService(ctx context.Context, userID uuid.UUID, name string) (Server, error)
	GetSrvUserService(ctx context.Context, userID uuid.UUID) ([]Server, error)
	GetInfoService(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Server, error)
	DeleteSrvService(ctx context.Context, id uuid.UUID, userID uuid.UUID) (Server, error)
}

type HandlerServer struct {
	service Service
}

func NewHandlerServer(service Service) *HandlerServer {
	return &HandlerServer{
		service: service,
	}
}

func (h *HandlerServer) CreateSrv(w http.ResponseWriter, r *http.Request) {
	var server Server
	ctx := r.Context()

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&server); err != nil {
		jwt.JsonError(w, "bad request", http.StatusBadRequest)
		return
	}

	userID, ok := ctx.Value("userID").(uuid.UUID)

	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.service.CreateSrvService(ctx, userID, server.Name)

	if err != nil {
		jwt.JsonError(w, "trouble with create server", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerServer) GetServerUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.service.GetSrvUserService(ctx, userID)

	if err != nil {
		jwt.JsonError(w, "trouble with get user servers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerServer) GetSrvInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := ctx.Value("userID").(uuid.UUID)

	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.service.GetInfoService(ctx, id, userID)

	if err != nil {
		jwt.JsonError(w, "trouble with getting server info", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *HandlerServer) DeleteServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := ctx.Value("userID").(uuid.UUID)

	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := h.service.DeleteSrvService(ctx, id, userID)

	if err != nil {
		jwt.JsonError(w, "trouble with getting server info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}
