package http_server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/bigflood/leaderboard/api"
	"github.com/bigflood/leaderboard/pkg/http_handler"
	"github.com/labstack/echo/v4"
)

type Server struct {
	httpServer *http.Server
	e          *echo.Echo
}

func New(lb api.LeaderBoard) *Server {
	handler := http_handler.New(lb)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	handler.Setup(e)

	httpServer := &http.Server{
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       330 * time.Second,
	}

	return &Server{
		httpServer: httpServer,
		e:          e,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	s.httpServer.Addr = addr
	return s.e.StartServer(s.httpServer)
}

func (s *Server) Serve(listener net.Listener) error {
	s.e.Listener = listener
	return s.e.StartServer(s.httpServer)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
