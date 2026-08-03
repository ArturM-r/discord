package channel

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID `json:"id"`
	Server_id uuid.UUID `json:"server_id"`
	Name      string    `json:"name"`
	CreatedAT time.Time `json:"created_at"`
}
