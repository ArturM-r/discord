package members

import (
	"context"
	"discord/internal/jwt"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type ServiceMbr interface {
	CreateMember(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) (Member, error)
	DeleteMember(ctx context.Context, userID uuid.UUID, targetID uuid.UUID, serverID uuid.UUID) (Member, error)
}

type Handler struct {
	service ServiceMbr
}

func NewHandlerMbr(service ServiceMbr) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(uuid.UUID)

	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := r.PathValue("id")

	if idStr == "" {
		jwt.JsonError(w, "id is empty", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)

	if err != nil {
		jwt.JsonError(w, "id is not valid", http.StatusBadRequest)
		return
	}

	response, err := h.service.CreateMember(ctx, userID, id)

	if err != nil {
		jwt.JsonError(w, "trouble with create", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		jwt.JsonError(w, "failed to encode", http.StatusInternalServerError)
		return
	}

}

func (h *Handler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	serverID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		jwt.JsonError(w, "invalid server id", http.StatusBadRequest)
		return
	}

	targetID, err := uuid.Parse(r.PathValue("uid"))
	if err != nil {
		jwt.JsonError(w, "invalid user id", http.StatusBadRequest)
		return
	}

	response, err := h.service.DeleteMember(r.Context(), userID, targetID, serverID)
	if err != nil {
		jwt.JsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
