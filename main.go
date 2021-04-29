package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bigflood/leaderboard/pkg/http_server"
	"github.com/bigflood/leaderboard/pkg/leaderboard"
)

func main() {
	const addr = ":8080"

	lb := &leaderboard.LeaderBoard{}

	server := http_server.New(lb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if waitSignal(ctx) {
			server.Shutdown(context.Background())
		}
	}()

	log.Println("listen", addr)
	err := server.ListenAndServe(addr)

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
