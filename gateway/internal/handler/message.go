package handler

import (
	"net/http"

	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) UpdateMessageStatus(c *gin.Context) {
	c.Status(200)
}

func (h *Handler) GetUnsentMessages(c *gin.Context) {
	c.Status(200)
}

// GetChatMessages godoc
// @Security BearerAuth
// @Summary Get Chat Messages
// @Description Get all messages for a specific chat
// @Tags messages
// @Accept json
// @Produce json
// @Param chat_id path string true "Chat ID"
// @Success 200 {object} models.GetUserMessageResponse
// @Failure 400 {object} errors.ErrorResponse "Bad Request"
// @Failure 401 {object} errors.ErrorResponse "Unauthorized"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /message/chat/{chat_id} [get]
func (h *Handler) GetChatMessages(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("GetChatMessages")

	chatID := c.Param("chat_id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
		return
	}

	res, err := h.service.Message.GetMessagesByChatID(chatID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error getting chat messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
