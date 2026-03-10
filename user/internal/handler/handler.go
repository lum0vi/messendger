package handler

import (
	"user/internal/middleware"
	"user/internal/service"

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
	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.GET("/me", h.GetMe)
	r.PUT("/me", h.UpdateMe)
	r.GET("/users", h.GetUsers)
	user := r.Group("/user")
	{
		user.POST("/id", h.GetUserById)
		user.POST("/name", h.GetUserByName)
	}

	return r
}
