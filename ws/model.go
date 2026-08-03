package ws

import (
	"discord/internal/message"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	UserID    uuid.UUID
	ChannelID uuid.UUID
	Conn      *websocket.Conn
	Send      chan []byte
}

type Hub struct {
	repo       *message.Repo
	channels   map[uuid.UUID]map[*Client]struct{}
	broadcast  chan message.Message
	register   chan *Client
	unregister chan *Client
}
