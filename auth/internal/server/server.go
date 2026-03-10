package server

import (
	"auth/internal/config"
	"auth/internal/handler"
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(cfg *config.Config, handler *handler.Handler) error {
	s.httpServer = &http.Server{
		Addr:           cfg.ServerHost + ":" + cfg.ServerPort,
		MaxHeaderBytes: 1 << 20,
		Handler:        handler.InitRoutes(),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
