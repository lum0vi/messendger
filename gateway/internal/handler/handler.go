package handler

import (
	"net/http"

	_ "gateway/docs" // тут будет документация
	"gateway/internal/config"
	"gateway/internal/kafka/producer"
	"gateway/internal/middleware"
	"gateway/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service     *service.Service
	upgrader    *websocket.Upgrader
	connections map[string]*websocket.Conn
	redisCon    *redis.Client
	cfg         *config.Config
	prod        *producer.ProducerKafka
}

func NewHandler(srv *service.Service, redisConn *redis.Client, cfg *config.Config, prod *producer.ProducerKafka) *Handler {
	return &Handler{
		service:     srv,
		connections: make(map[string]*websocket.Conn),
		redisCon:    redisConn,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // разрешаем все соединения, для разработки ок
			},
		},
		cfg:  cfg,
		prod: prod,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestMiddleware())
	r.GET("/ws", middleware.AuthMiddleware(), h.Websocket)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// auth := r.Group("/auth")
	// {
	// 	auth.POST("/register", h.Register)
	// 	auth.POST("/login", h.Login)
	// }
	// user := r.Group("/user")
	// {
	// 	user.Use(middleware.AuthMiddleware())
	// 	user.GET("/me", h.GetMe)
	// 	user.PUT("/me", h.UpdateMe)
	// 	user.GET("/users", h.GetUsers)
	// 	user.POST("/id", h.GetUserByID)
	// 	user.POST("/name", h.GetUserByUsername)
	// }

	// chat := r.Group("/chat")
	// {
	// 	chat.Use(middleware.AuthMiddleware())
	// 	chat.GET("/", h.GetMeChats)
	// 	chat.POST("/private", h.CreatePrivateChat)
	// 	chat.POST("/public", h.CreatePublicChat)
	// 	chat.GET("/:chat_id/users", h.GetChatUsers)
	// }

	// message := r.Group("/message")
	// {
	// 	message.Use(middleware.AuthMiddleware())
	// 	message.PUT("/upd", h.UpdateMessageStatus)
	// 	message.GET("/:id", h.GetUnsentMessages)
	// 	message.GET("/chat/:chat_id", h.GetChatMessages)
	// }

	return r
}
