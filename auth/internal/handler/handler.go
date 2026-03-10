package handler

import (
	"auth/internal/middleware"
	"auth/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
	return r
}
