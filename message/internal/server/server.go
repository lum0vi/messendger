package server

import (
	"context"
	"message/internal/config"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(hand http.Handler, cfg *config.Config) error {
	s.httpServer = &http.Server{
		Addr:           cfg.Host + ":" + cfg.Port,
		Handler:        hand,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
