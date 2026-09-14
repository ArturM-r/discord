package ws

import (
	"context"
	"discord/internal/checkmember"
	"discord/internal/message"
	"encoding/json"

	"github.com/google/uuid"
)

func NewHub(repo *message.Repo, memberCache *checkmember.MemberCache) *Hub {
	return &Hub{
		repo:        repo,
		servers:     make(map[uuid.UUID]map[*Client]struct{}),
		broadcast:   make(chan BroadcastMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		membercache: memberCache,
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, clients := range h.servers {
				for client := range clients {
					close(client.Send)
					client.Cancel()
				}
			}
			return

		case client := <-h.register:
			if h.servers[client.ServerID] == nil {
				h.servers[client.ServerID] = make(map[*Client]struct{})
			}
			h.servers[client.ServerID][client] = struct{}{}

		case client := <-h.unregister:
			if clients, ok := h.servers[client.ServerID]; ok {
				delete(clients, client)
				close(client.Send)

				if len(clients) == 0 {
					delete(h.servers, client.ServerID)
				}
			}
			client.Cancel()

		case bm := <-h.broadcast:

			data, err := json.Marshal(bm.Msg)
			if err != nil {
				continue
			}

			for client := range h.servers[bm.ServerID] {
				select {
				case client.Send <- data:
				default:
					close(client.Send)
					delete(h.servers[client.ServerID], client)
					client.Cancel()
				}
			}
		}
	}

}
