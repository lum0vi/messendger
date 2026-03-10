package handler

import (
	customErr "gateway/internal/errors"
	"gateway/internal/middleware"
	"gateway/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CreatePrivateChat godoc
// @Security BearerAuth
// @Summary Create Private Chat
// @Description Create a private chat between users
// @Tags chats
// @Accept json
// @Produce json
// @Param input body models.CreatePrivateChatRequest true "Private Chat Data"
// @Success 201 {object} models.CreatePrivateChatResponse
// @Failure 400 {object} errors.ErrorResponse "Bad Request"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /chat/private [post]
func (h *Handler) CreatePrivateChat(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("CreatePrivateChat")
	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetMe: failed to get user id")
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	var req models.CreatePrivateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("CreatePrivateChat: ShouldBindJSON: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusBadRequest, "invalid request data")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	res, err := h.service.Chat.CreatePrivateChat(id.(string), &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetMe: CreatePrivateChat error: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Infof("handler.GetMe: CreatePrivateChat: %v", res)
	c.JSON(http.StatusCreated, res)
}

// CreatePublicChat godoc
// @Security BearerAuth
// @Summary Create Public Chat
// @Description Create a public group chat
// @Tags chats
// @Accept json
// @Produce json
// @Param input body models.CreatePublicChatRequest true "Public Chat Data"
// @Success 201 {object} models.CreatePrivateChatResponse
// @Failure 400 {object} errors.ErrorResponse "Bad Request"
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /chat/public [post]
func (h *Handler) CreatePublicChat(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Infoln("CreatePublicChat")
	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetMe: failed to get user id")
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	var req models.CreatePublicChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("CreatePublicChat: ShouldBindJSON: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusBadRequest, "invalid request data")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	res, err := h.service.Chat.CreatePublicChat(id.(string), &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetMe: CreatePublicChat error: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Infof("handler.GetMe: CreatePublicChat: %v", res)
	c.JSON(http.StatusCreated, res)
}

// GetMeChats godoc
// @Security BearerAuth
// @Summary Get My Chats
// @Description Retrieve all user chats
// @Tags chats
// @Accept json
// @Produce json
// @Success 200 {object} models.GetMeChatsResponse
// @Failure 500 {object} errors.ErrorResponse "Internal Server Error"
// @Router /chat [get]
func (h *Handler) GetMeChats(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetMeChats")
	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetMeChats: failed to get user id")
		ErrResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, ErrResp)
		return
	}
	var req models.GetMeChatsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetMeChats: ShouldBindQuery: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	res, err := h.service.Chat.GetMeChats(id.(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetMeChats: GetMeChats error: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetChatUsers godoc
// @Security BearerAuth
// @Summary      GetChatUsers
// @Description  получить пользователей чата
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param chat_id path string true "ID чата"
// @Success      200  {object}  models.GetChatUsersResponse
// @Failure default {object} errors.ErrorResponse
// @Router       /chat/{chat_id}/users [get]
func (h *Handler) GetChatUsers(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetChatUsers")

	chatID := c.Param("chat_id")

	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetChatUsers: failed to get user id")
		ErrResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, ErrResp)
		return
	}
	var req models.GetChatUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetChatUsers: ShouldBindQuery: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	res, err := h.service.Chat.GetChatUsers(id.(string), chatID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("handler.GetChatUsers: GetChatUsers error: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	c.JSON(http.StatusOK, res)
}
