package main

import (
	"context"
	"github.com/bigflood/leaderboard/pkg/http_server"
	"github.com/bigflood/leaderboard/pkg/leaderboard"
	"github.com/bigflood/leaderboard/pkg/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const addr = ":8080"

	lb := &leaderboard.LeaderBoard{
		Storage: &storage.MemStorage{},
	}

	server := http_server.New(lb, log.Default())

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
