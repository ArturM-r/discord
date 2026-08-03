package ws

import (
	"context"
	"discord/internal/jwt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) WsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)

	if !ok {
		jwt.JsonError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	channelID, err := uuid.Parse(r.URL.Query().Get("channel_id"))
	if err != nil {
		http.Error(w, "invalid channel_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	client := &Client{
		UserID:    userID,
		ChannelID: channelID,
		Conn:      conn,
		Send:      make(chan []byte, 256),
	}

	h.register <- client

	go client.WritePump()

	go client.ReadPump(r.Context(), h)

}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}

func (c *Client) ReadPump(ctx context.Context, hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()

	go func() {
		<-ctx.Done()
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		msg, err := hub.repo.WriteMessage(ctx, c.ChannelID, c.UserID, string(data))
		if err != nil {
			continue
		}

		hub.broadcast <- msg
	}
}
