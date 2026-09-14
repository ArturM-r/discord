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

	mux, err := initroute.Init(db, cfg.HMACKey, ongoingCtx)
	if err != nil {
		log.Fatal(err)
	}

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
	// server.Shutdown waits for in-flight requests to finish, and those
	// requests run under ongoingCtx (it's the server's BaseContext and the
	// DB pool's context). Cancelling it before Shutdown returns - as this
	// used to do - killed every in-flight request and DB query immediately
	// instead of giving them the shutdown window to complete.
	err = server.Shutdown(shutdownCtx)

	if err != nil {
		log.Println("Forced shutdown after timeout")
		time.Sleep(_shutdownHardPeriod)
	}

	stopOngoingGracefully()

	log.Println("Bye!")
}
