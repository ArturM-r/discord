package server

import (
	"time"

	"github.com/google/uuid"
)

type Server struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAT time.Time `json:"created_at"`
}
