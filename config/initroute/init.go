package initroute

import (
	"context"
	"discord/internal/channel"
	"discord/internal/checkmember"
	"discord/internal/jwt"
	"discord/internal/members"
	"discord/internal/message"
	"discord/internal/server"
	"discord/internal/user"
	"discord/ws"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(db *pgxpool.Pool, secret string, ctx context.Context) (http.Handler, error) {
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db), secret))
	channelHandler := channel.NewHandler(channel.NewService(channel.NewChannelPool(db), db))
	messageHandler := message.NewHandler(message.NewService(message.NewMessagePool(db)))

	msgRepo := message.NewMessagePool(db)

	m, err := checkmember.Unload(ctx, db)

	if err != nil {
		return nil, fmt.Errorf("memberchache unload problem: %w", err)
	}

	MemberCache := checkmember.NewMemberCache(m)
	serverHandler := server.NewHandlerServer(server.NewService(server.NewMessagePool(db), db, MemberCache))
	memberHandler := members.NewHandlerMbr(members.NewServiceMbr(members.NewRepository(db), db, MemberCache))

	hub := ws.NewHub(msgRepo, MemberCache)
	go hub.Run(ctx)

	mux := http.NewServeMux()

	mux.HandleFunc("/auth/register", userHandler.Register)
	mux.HandleFunc("/auth/login", userHandler.Login)

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

	mux.HandleFunc("/channels/{id}/messages", messageHandler.GetMessages)

	mux.HandleFunc("/ws", hub.WsHandler)

	authMiddleware := jwt.NewAuthMiddleware(secret)
	return authMiddleware.Authorize(mux), nil
}
