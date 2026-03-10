package handler

import (
	_ "chat/docs"
	"chat/internal/middleware"
	"chat/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(middleware.LoggerMiddleware())
	chats := router.Group("/chat")
	{
		chats.POST("/private", h.CreatePrivateChat)
		chats.POST("/public", h.CreatePublicChat)
		chats.GET("/", h.GetChats)
		chats.GET("/:chat_id/users", h.GetChatUsers)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
