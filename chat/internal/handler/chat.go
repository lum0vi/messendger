package handler

import (
	"chat/internal/errors"
	"chat/internal/middleware"
	"chat/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @Security ApiKeyAuth
// @Summary Создание Приватного чата
// @Tags API создание чата
// @Description Создание приватного чата
// @ID create-private-chat
// @Accept  json
// @Produce  json
// @Param input body models.CreatePrivateChatRequest true "credentials"
// @Success 200 {object} models.CreatePrivateChatResponse "data"
// @Failure default {object} errors.ErrorResponse
// @Router /chat/private [post]
func (h *Handler) CreatePrivateChat(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("CreateChat")
	id := c.GetHeader("id")
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"request_id": requestId,
		}).Error("user ID required")
		errResp := errors.NewErrorResponse(http.StatusUnauthorized, "id is required")
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)

	var req *models.CreatePrivateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("invalid request body, %v", err)
		errResp := errors.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Infof("request body: %v", req)
	res, err := h.svc.Chat.CreatePrivateChat(id, req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("create chat failed, %v", err)
		errCode, msg := errors.ParseCustomError(err)
		errResp := errors.NewErrorResponse(errCode, msg)
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Infof("CreateChat response %+v", res)
	var resp models.CreatePrivateChatResponse
	resp.ChatID = res

	c.JSON(http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// @Summary Создание Группы
// @Tags API создание чата
// @Description Создание публичного чата
// @ID create-public-chat
// @Accept  json
// @Produce  json
// @Param input body models.CreatePublicChatRequest true "credentials"
// @Success 200 {object} models.CreatePublicChatResponse "data"
// @Failure default {object} errors.ErrorResponse
// @Router /chat/public [post]
func (h *Handler) CreatePublicChat(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Info("CreatePublicChat")
	id := c.GetHeader("id")
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("user ID required")
		errResp := errors.NewErrorResponse(http.StatusUnauthorized, "id is required")
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)
	var req *models.CreatePublicChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("invalid request body, %v", err)
		errResp := errors.NewErrorResponse(http.StatusBadRequest, "invalid request body")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Infof("request body: %+v", req)
	res, err := h.svc.Chat.CreatePublicChat(id, req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("create chat failed, %v", err)
		errResp := errors.NewErrorResponse(http.StatusInternalServerError, "create chat failed")
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Infof("CreateChat response %+v", res)
	var resp models.CreatePublicChatResponse
	resp.ChatID = res
	c.JSON(http.StatusCreated, resp)
}

// @Security ApiKeyAuth
// @Summary Получить чаты пользователя
// @Tags API получить чаты
// @Description Получение чатов пользователя
// @ID get-chats
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GetChatsResponse "data"
// @Failure default {object} errors.ErrorResponse
// @Router /chat [get]
func (h *Handler) GetChats(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("GetChats")
	id := c.GetHeader("id")
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Error("user ID required")
		errResp := errors.NewErrorResponse(http.StatusUnauthorized, "id is required")
		c.JSON(http.StatusUnauthorized, errResp)
		return
	}
	res, err := h.svc.Chat.GetChats(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("get chat failed, %v", err)
		errResp := errors.NewErrorResponse(http.StatusInternalServerError, "get chat failed")
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	var resp models.GetChatsResponse
	resp.ChatID = res
	c.JSON(http.StatusOK, resp)
}

// GetChatUsers получает список участников чата по chat_id
// @Summary Получить участников чата
// @Tags chats
// @Param chat_id path string true "ID чата"
// @Success 200 {array} models.GetChatUsersResponse
// @Failure default {object} errors.ErrorResponse
// @Router /chat/{chat_id}/users [get]
func (h *Handler) GetChatUsers(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("GetChatUsers")
	chatID := c.Param("chat_id")
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("chat_id: " + chatID)
	res, err := h.svc.Chat.GetUsersChat(chatID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId": requestId,
		}).Errorf("get chat users failed, %v", err)
		errResp := errors.NewErrorResponse(http.StatusInternalServerError, "get chat users failed")
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestId": requestId,
	}).Infof("GetChatUsers response %+v", res)
	var resp models.GetChatUsersResponse
	resp.Users = res
	c.JSON(http.StatusOK, resp)
}
