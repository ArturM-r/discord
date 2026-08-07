package initroute

import (
	"context"
	"discord/internal/channel"
	"discord/internal/checkmember"
	"discord/internal/members"
	"discord/internal/message"
	"discord/internal/server"
	"discord/internal/user"
	"discord/ws"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(db *pgxpool.Pool, secret string, ctx context.Context) (*http.ServeMux, error) {
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db), secret))
	serverHandler := server.NewHandlerServer(server.NewService(server.NewMessagePool(db)))
	memberHandler := members.NewHandlerMbr(members.NewServiceMbr(members.NewRepository(db)))
	channelHandler := channel.NewHandler(channel.NewService(channel.NewChannelPool(db)))
	messageHandler := message.NewHandler(message.NewService(message.NewMessagePool(db)))

	msgRepo := message.NewMessagePool(db)

	m, err := checkmember.Unload(ctx, db)

	MemberCache := checkmember.NewMemberCache(m)

	if err != nil {
		return nil, fmt.Errorf("memberchache unload problem: %w", err)
	}

	hub := ws.NewHub(msgRepo, MemberCache)
	go hub.Run(ctx)

	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("/auth/register", userHandler.Register)
	mux.HandleFunc("/auth/login", userHandler.Login)

	// servers
	mux.HandleFunc("/servers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			serverHandler.GetServerUser(w, r)
		case http.MethodPost:
			serverHandler.CreateSrv(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	mux.HandleFunc("/servers/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			serverHandler.GetSrvInfo(w, r)
		case http.MethodDelete:
			serverHandler.DeleteServer(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})

	// members
	mux.HandleFunc("/servers/{id}/members", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			memberHandler.Create(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	mux.HandleFunc("/servers/{id}/members/{uid}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			memberHandler.DeleteMember(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})

	// channels
	mux.HandleFunc("/servers/{id}/channels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			channelHandler.GetChannel(w, r)
		case http.MethodPost:
			channelHandler.CreateChannel(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	mux.HandleFunc("/servers/{id}/channels/{cid}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			channelHandler.DeleteChannel(w, r)
		default:
			http.Error(w, "Method Not Allowed", 405)
		}
	})

	// messages
	mux.HandleFunc("/channels/{id}/messages", messageHandler.GetMessages)

	// websocket
	mux.HandleFunc("/ws", hub.WsHandler)

	return mux, nil
}
