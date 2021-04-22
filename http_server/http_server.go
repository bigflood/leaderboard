package http_server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bigflood/leaderboard/http_handler"
	"github.com/bigflood/leaderboard/leaderboard"
)

type Server struct {
	httpServer *http.Server
}

func New(lb *leaderboard.LeaderBoard) *Server {
	handler := http_handler.New(lb)
	httpServer := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       330 * time.Second,
		Handler:           handler.Setup(),
	}

	return &Server{
		httpServer: httpServer,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	s.httpServer.Addr = addr
	return s.httpServer.ListenAndServe()
}

func (s *Server) Serve(listener net.Listener) error {
	return s.httpServer.Serve(listener)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
