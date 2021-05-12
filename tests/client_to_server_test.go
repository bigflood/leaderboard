package tests

import (
	"context"
	"github.com/bigflood/leaderboard/api"
	"github.com/bigflood/leaderboard/pkg/http_client"
	"github.com/bigflood/leaderboard/pkg/http_server"
	"github.com/bigflood/leaderboard/pkg/leaderboard"
	"github.com/bigflood/leaderboard/pkg/storage"
	"log"
	"net"
	"testing"
	"time"
)

func TestClientToServerLeaderBoard(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)

	lb := &leaderboard.LeaderBoard{
		NowFunc: func() time.Time {
			return now
		},
		Storage: &storage.MemStorage{},
	}

	testClientToServer(t, lb, func(client api.LeaderBoard) {
		testLeaderBoard(t, client, now)
	})
}

func TestClientToServerLeaderBoardWithMultiGoroutines(t *testing.T) {
	lb := &leaderboard.LeaderBoard{
		Storage: &storage.MemStorage{},
	}

	testClientToServer(t, lb, func(client api.LeaderBoard) {
		testMultiGoroutines(t, client)
	})
}

func testClientToServer(t *testing.T, logic api.LeaderBoard, f func(client api.LeaderBoard)) {
	server := http_server.New(logic, nil)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	go server.Serve(listener)
	defer server.Shutdown(context.Background())

	client := http_client.New("http://" + listener.Addr().String())

	f(client)
}
