package tests

import (
	// "context"
	// "log"
	// "net"

	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/bigflood/leaderboard/api"
	"github.com/bigflood/leaderboard/pkg/http_client"
	"github.com/bigflood/leaderboard/pkg/http_server"
	"github.com/bigflood/leaderboard/pkg/leaderboard"
)

func TestClientToServerLeaderBoard(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)

	lb := &leaderboard.LeaderBoard{
		NowFunc: func() time.Time {
			return now
		},
	}

	testClientToServer(t, lb, func(client api.LeaderBoard) {
		testLeaderBoard(t, client, now)
	})
}

func TestClientToServerLeaderBoardWithMultiGoroutines(t *testing.T) {
	lb := &leaderboard.LeaderBoard{}

	testClientToServer(t, lb, func(client api.LeaderBoard) {
		testMultiGoroutines(t, client)
	})
}

func testClientToServer(t *testing.T, logic api.LeaderBoard, f func(client api.LeaderBoard)) {
	server := http_server.New(logic)

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
