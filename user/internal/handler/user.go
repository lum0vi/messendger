package handler

import (
	"net/http"
	"user/internal/errors"
	"user/internal/middleware"
	"user/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetMe(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	id := c.GetHeader("id")
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"request_id": requestId,
		}).Info("Request ID required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id is required"})
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)

	res, err := h.svc.User.GetMe(id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateMe(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	id := c.GetHeader("id")
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"request_id": requestId,
		}).Info("Request ID required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id is required"})
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)
	var req models.UpdateMeRequest
	if err := c.BindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	err := h.svc.User.UpdateMe(id, &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *Handler) GetUsers(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	id := c.GetHeader("id")
	if id == "" {
		logrus.WithFields(logrus.Fields{
			"request_id": requestId,
		}).Info("Request ID required")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "id is required"})
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("id: " + id)
	res, err := h.svc.User.GetUsers()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Error("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUserById(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("GetUserById")
	var req models.GetUserByIDRequest
	if err := c.BindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
	}
	res, err := h.svc.User.GetUserByID(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetUserByName(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestID)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestID: requestId,
	}).Info("GetUserById")
	var req models.GetUserByUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
	}
	res, err := h.svc.User.GetUserByUsername(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestID: requestId,
		}).Info("err: " + err.Error())
		code, msg := errors.ParseCustomError(err)
		c.JSON(code, gin.H{"error": msg})
	}
	c.JSON(http.StatusOK, res)
}
