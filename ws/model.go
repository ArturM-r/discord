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

// BroadcastMessage carries the server a message belongs to alongside the
// message itself, since h.servers is keyed by server ID and message.Message
// only carries a channel ID - broadcasting by msg.ChannelID against a
// server-keyed map never matched any connected client.
type BroadcastMessage struct {
	ServerID uuid.UUID
	Msg      message.Message
}

// incomingMessage is the JSON payload a client sends over the websocket
// connection. A channel ID is required because messages.channel_id is a
// NOT NULL foreign key - without it, every inbound message failed to persist.
type incomingMessage struct {
	ChannelID uuid.UUID `json:"channel_id"`
	Content   string    `json:"content"`
}
