package message

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	ChannelID uuid.UUID `json:"channel_id"`
	UserID    uuid.UUID `json:"user_id"`
	UserName  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAT time.Time `json:"created_at"`
}
