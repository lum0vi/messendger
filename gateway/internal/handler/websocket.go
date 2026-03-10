package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Websocket godoc
// @Security BearerAuth
// @Summary      WebSocket
// @Description  Establish a websocket connection for real-time messaging.
// @Tags         websocket
// @Accept       json
// @Produce      json
// @Success      101 {object} string "Connection Established"
// @Failure      400 {object} string "Bad Request"
// @Failure      401 {object} string "Unauthorized"
// @Router       /ws [get]
func (h *Handler) Websocket(c *gin.Context) {
	requestID, exists := c.Get(middleware.RequestIDKey)
	if !exists {
		requestID = "unknown"
	}

	// Получаем userID из контекста после auth middleware
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Error("userID not found in context")
		c.JSON(400, gin.H{"error": "userID not found"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Error("userID is not a string")
		c.JSON(400, gin.H{"error": "invalid userID"})
		return
	}

	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestID,
	}).Infof("Handle Websocket, userID: %v", userIDStr)
	socketID := uuid.NewString()

	h.redisCon.Set(c, "socket:"+userIDStr, socketID, time.Minute*30)
	fmt.Println(userID)
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestID,
	}).Info("Save redis socket")

	// Сохраняем соответствие socketID -> userID для поиска при отправке
	h.redisCon.Set(c, "socket_to_user:"+socketID, userIDStr, time.Minute*30)

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Error(fmt.Sprintf("error upgrading connection: %v", err))
		c.JSON(400, gin.H{"error": "failed to upgrade connection"})
		return
	}
	defer conn.Close()
	h.connections[socketID] = conn

	logrus.WithFields(logrus.Fields{
		"socketID":              socketID,
		middleware.RequestIDKey: requestID,
	}).Info(fmt.Sprintf("User %s connected with socketID %s\n", userIDStr, socketID))

	for {
		if h.redisCon.Get(c, "socket:"+userIDStr).Err() != nil {
			logrus.WithFields(logrus.Fields{
				middleware.RequestIDKey: requestID,
			}).Info("connection closed")
			break
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"socketID":              socketID,
				middleware.RequestIDKey: requestID,
			}).Error("ReadMessage error:", err)
			break
		}
		logrus.WithFields(logrus.Fields{
			"socketID":              socketID,
			middleware.RequestIDKey: requestID,
		}).Info(string(msg))

		if err := conn.WriteMessage(websocket.PongMessage, []byte("Message received")); err != nil {
			log.Println("Error sending message:", err)
			break
		}

		// Парсим сообщение с учетом структуры {"type": "send_message", "data": {...}}
		var rawMsg struct {
			Type string `json:"type"`
			Data struct {
				ChatID   string `json:"chat_id"`
				SenderID string `json:"sender_id"`
				Content  string `json:"content"`
				SentAt   int64  `json:"sent_at"`
			} `json:"data"`
		}

		if err := json.Unmarshal(msg, &rawMsg); err != nil {
			logrus.Info("Error json decoding message:", err)
			break
		}

		// Отправляем исходное сообщение в Kafka
		if err := h.prod.Send(c, []byte(rawMsg.Data.ChatID), msg); err != nil {
			logrus.WithFields(logrus.Fields{
				middleware.RequestIDKey: requestID,
			}).Error(fmt.Sprintf("error sending message: %v", err))
			break
		}
	}
	delete(h.connections, socketID)
	h.redisCon.Del(c, "socket:"+userIDStr)
	h.redisCon.Del(c, "socket_to_user:"+socketID)
}
