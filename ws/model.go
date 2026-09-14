package ws

import (
	"context"
	"discord/internal/checkmember"
	"discord/internal/message"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID   uuid.UUID
	ServerID uuid.UUID
	Conn     *websocket.Conn
	Send     chan []byte
	Cancel   context.CancelFunc
}

type Hub struct {
	repo        *message.Repo
	servers     map[uuid.UUID]map[*Client]struct{}
	broadcast   chan BroadcastMessage
	register    chan *Client
	unregister  chan *Client
	membercache *checkmember.MemberCache
}

type BroadcastMessage struct {
	ServerID uuid.UUID
	Msg      message.Message
}

type incomingMessage struct {
	ChannelID uuid.UUID `json:"channel_id"`
	Content   string    `json:"content"`
}
