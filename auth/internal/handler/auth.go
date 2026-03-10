package handler

import (
	"auth/internal/errors"
	"net/http"

	"auth/internal/middleware"
	"auth/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) Register(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestIdKey)
	if !ok {
		requestId = "Unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIdKey: requestId,
	}).Info("Register")
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIdKey: requestId,
		}).Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.srv.Auth.Register(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIdKey: requestId,
		}).Errorf("Failed to register user: %v", err)
		code, msg, cusomErr := errors.ParseHttpError(err)
		if cusomErr == nil {
			c.JSON(code, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIdKey: requestId,
	}).Infoln("Successfully registered")
	c.JSON(200, gin.H{"id": id})
}

func (h *Handler) Login(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestIdKey)
	if !ok {
		requestId = "Unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIdKey: requestId,
	}).Info("Login")
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIdKey: requestId,
		}).Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.srv.Auth.Login(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIdKey: requestId,
		}).Errorf("Failed to login: %+v", err)
		code, msg, customErr := errors.ParseHttpError(err)
		if customErr == nil {
			c.JSON(code, gin.H{"error": msg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIdKey: requestId,
	}).Infof("Successfully login")
	c.JSON(http.StatusOK, gin.H{"token": token})
}
