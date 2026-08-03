package ws

import (
	"context"
	"discord/internal/message"
	"encoding/json"

	"github.com/google/uuid"
)

func NewHub(repo *message.Repo) *Hub {
	return &Hub{
		repo:       repo,
		channels:   make(map[uuid.UUID]map[*Client]struct{}),
		broadcast:  make(chan message.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, clients := range h.channels {
				for client := range clients {
					close(client.Send)
				}
			}
			return

		case client := <-h.register:
			if h.channels[client.ChannelID] == nil {
				h.channels[client.ChannelID] = make(map[*Client]struct{})
			}
			h.channels[client.ChannelID][client] = struct{}{}

		case client := <-h.unregister:
			if clients, ok := h.channels[client.ChannelID]; ok {
				delete(clients, client)
				close(client.Send)

				if len(clients) == 0 {
					delete(h.channels, client.ChannelID)
				}
			}

		case msg := <-h.broadcast:

			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			for client := range h.channels[msg.ChannelID] {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.channels[client.ChannelID], client)
				}
			}
		}
	}

}
