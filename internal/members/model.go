package members

import "github.com/google/uuid"

type Member struct {
	ID       uuid.UUID `json:"id"`
	ServerID uuid.UUID `json:"server_id"`
	UserID   uuid.UUID `json:"user_id"`
	Role     string    `json:"role"` //owner, admin, member
}
