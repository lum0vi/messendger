package handler

import (
	"message/internal/middleware"
	"message/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.RequestMiddleware())
	mess := router.Group("/message")
	{
		mess.GET("/user/:user_id", h.GetUserMessage)
		mess.GET("/chat/:chat_id", h.GetMessagesByChatID)
	}

	return router
}
