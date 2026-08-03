package main

import (
	"context"
	"discord/config"
	"discord/config/initroute"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	db := config.DbConn(ongoingCtx)
	cfg := config.GetConfig()
	config.RunMigrations(cfg.DatabaseUrl)

	mux := initroute.Init(db, cfg.HMACKey, ongoingCtx)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Server error:", err)
			os.Exit(1)
		}
	}()

	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Shutting down...")

	time.Sleep(_readinessDrainDelay)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	stopOngoingGracefully()
	err := server.Shutdown(shutdownCtx)

	if err != nil {
		log.Println("Forced shutdown after timeout")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Bye!")
}
