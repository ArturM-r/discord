package message

import (
	"discord/internal/jwt"
	"encoding/json"

	"context"
	"net/http"

	"github.com/google/uuid"
)

type Service interface {
	GetByChannel(ctx context.Context, channelID string, userID uuid.UUID) ([]Message, error)
}

type HandlerMSG struct {
	service Service
}

func NewHandler(service Service) *HandlerMSG {
	return &HandlerMSG{
		service: service,
	}
}

func (m *HandlerMSG) GetMessages(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()

	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if id == "" {
		jwt.JsonError(w, "id is empty", http.StatusBadRequest)
		return
	}

	messages, err := m.service.GetByChannel(ctx, id, userID)
	if err != nil {
		jwt.JsonError(w, "trouble with messages taken", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(messages); err != nil {
		jwt.JsonError(w, "failed to endcode", http.StatusInternalServerError)
		return
	}
}
