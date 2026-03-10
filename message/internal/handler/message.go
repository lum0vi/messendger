package handler

import (
	"message/internal/errors"
	"message/internal/middleware"
	"message/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetMessagesByChatID(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithField(middleware.RequestIDKey, requestID).Info("GetMessagesByChatID")
	chatID := c.Param("chat_id")
	if len(chatID) == 0 {
		errResp := errors.NewErrorResponse(http.StatusBadRequest, "invalid request")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	res, err := h.svc.Message.GetMessagesByChatID(chatID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		})
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Errorf("error getting messages by chatID: %v", err)
		errResp := errors.NewErrorResponse(errors.ParseCustomError(err))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	resp := models.GetUserMessageResponse{
		Messages: res,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetUserMessage(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithField(middleware.RequestIDKey, requestID).Info("GetUserMessage")
	userID := c.Param("user_id")
	if len(userID) == 0 {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Errorf("invalid url")
		errResp := errors.NewErrorResponse(http.StatusBadRequest, "invalid request")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}

	res, err := h.svc.Message.GetUserMessages(userID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestID,
		}).Errorf("error getting messages by userID: %v", err)
		errResp := errors.NewErrorResponse(errors.ParseCustomError(err))
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	c.JSON(http.StatusOK, res)
}
