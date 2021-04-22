package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bigflood/leaderboard/http_handler"
	"github.com/bigflood/leaderboard/leaderboard"
)

func main() {
	lb := &leaderboard.LeaderBoard{}
	handler := http_handler.New(lb)
	server := http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       330 * time.Second,
		Handler:           handler.Setup(),
	}

	log.Println("listen", server.Addr)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if waitSignal(ctx) {
			server.Shutdown(context.Background())
		}
	}()

	err := server.ListenAndServe()

	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func waitSignal(ctx context.Context) bool {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		log.Println("signal", sig)
		return true
	case <-ctx.Done():
		return false
	}
}
