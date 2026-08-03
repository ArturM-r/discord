package channel

import (
	"context"
	"discord/internal/jwt"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	GetChService(ctx context.Context, serverID uuid.UUID) ([]Channel, error)
	CreateChService(ctx context.Context, serverID uuid.UUID, name string) (Channel, error)
	DeleteChnService(ctx context.Context, serverID uuid.UUID, channelID uuid.UUID) (Channel, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetChannel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := uuid.Parse(idStr)

	if err != nil {
		jwt.JsonError(w, "not valid id", http.StatusBadRequest)
		return
	}

	response, err := h.service.GetChService(r.Context(), id)

	if err != nil {
		jwt.JsonError(w, "trouble with channel taken", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to endcode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	serverIDStr := r.PathValue("id")
	serverID, err := uuid.Parse(serverIDStr)
	if err != nil {
		jwt.JsonError(w, "not valid id", http.StatusBadRequest)
		return
	}

	channelIDStr := r.PathValue("cid")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		jwt.JsonError(w, "not valid id", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteChnService(r.Context(), serverID, channelID)
	if err != nil {
		jwt.JsonError(w, "trouble with channel delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	serverIDStr := r.PathValue("id")
	serverID, err := uuid.Parse(serverIDStr)
	if err != nil {
		jwt.JsonError(w, "not valid id", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var channel Channel
	if err := json.NewDecoder(r.Body).Decode(&channel); err != nil {
		jwt.JsonError(w, "bad request", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(channel.Name) == "" {
		jwt.JsonError(w, "channel name is required", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateChService(r.Context(), serverID, channel.Name)
	if err != nil {
		jwt.JsonError(w, "failed to create channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}
}
