package handler

import (
	"net/http"

	customErr "gateway/internal/errors"
	"gateway/internal/middleware"
	"gateway/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetMe godoc
// @Security BearerAuth
// @Summary      Get User Information
// @Description  Retrieve the currently authenticated user's data.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.LoginResponse
// @Failure      401  {object}  errors.ErrorResponse "Unauthorized"
// @Failure      500  {object}  errors.ErrorResponse "Internal Server Error"
// @Router       /user/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetMe")

	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetMe: failed to get user id")
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}

	res, err := h.service.User.GetMe(id.(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error getting user: %v", err)
		code, msg := customErr.ParseCustomError(err)
		errResp := customErr.NewErrorResponse(code, msg)
		c.JSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateMe godoc
// @Security BearerAuth
// @Summary      Update User Information
// @Description  Update the currently authenticated user's data.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body  models.UpdateMeRequest  true  "User Update Data"
// @Failure      400  {object}  errors.ErrorResponse "Bad Request"
// @Failure      401  {object}  errors.ErrorResponse "Unauthorized"
// @Failure      500  {object}  errors.ErrorResponse "Internal Server Error"
// @Router       /user/me [put]
func (h *Handler) UpdateMe(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.UpdateMe")
	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.UpdateMe: failed to get user id")
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	var req models.UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error parsing request: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusBadRequest, "bad request")
		c.JSON(http.StatusOK, errResp)
		return
	}
	err := h.service.User.UpdateMe(id.(string), &req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error updating user: %v", err)
		code, msg := customErr.ParseCustomError(err)
		errResp := customErr.NewErrorResponse(code, msg)
		c.JSON(http.StatusOK, errResp)
		return
	}
	c.Status(http.StatusOK)
}

// GetUsers godoc
// @Security BearerAuth
// @Summary      Get All Users
// @Description  Retrieve a list of all users.
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.GetUsersResponse
// @Failure      401  {object}  errors.ErrorResponse "Unauthorized"
// @Failure      500  {object}  errors.ErrorResponse "Internal Server Error"
// @Router       /user/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetUsers")
	id, ok := c.Get(middleware.UserIDKey)
	if !ok {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Error("handler.GetUsers: failed to get user id")
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	res, err := h.service.User.GetUsers(id.(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error getting users: %v", err)
		code, msg := customErr.ParseCustomError(err)
		errResp := customErr.NewErrorResponse(code, msg)
		c.JSON(http.StatusOK, errResp)
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetUserByID godoc
// @Security BearerAuth
// @Summary      Get User By ID
// @Description  Retrieve user data by user ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body  models.GetUserByIdRequest  true  "id пользователя"
// @Success      200  {object}  models.GetUserByIdResponse
// @Failure default {object} errors.ErrorResponse
// @Router       /user/id [post]
func (h *Handler) GetUserByID(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetUserByID")
	var req models.GetUserByIdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error parsing request: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	res, err := h.service.User.GetUserById(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		})
		code, msg := customErr.ParseCustomError(err)
		errResp := customErr.NewErrorResponse(code, msg)
		c.JSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetUserByUsername godoc
// @Security BearerAuth
// @Summary      GetUserByUsername
// @Description  получить пользователя по username
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input  body  models.GetUserByUsernameRequest  true  "username пользователя"
// @Success      200  {object}  models.GetUserByUsernameResponse
// @Failure default {object} errors.ErrorResponse
// @Router       /user/name [post]
func (h *Handler) GetUserByUsername(c *gin.Context) {
	requestID, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestID = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		"requestID": requestID,
	}).Info("handler.GetUserByUsername")
	var req models.GetUserByUsernameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error parsing request: %v", err)
		errResp := customErr.NewErrorResponse(http.StatusInternalServerError, "internal server error")
		c.JSON(http.StatusOK, errResp)
		return
	}
	res, err := h.service.User.GetUserByUsername(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestID": requestID,
		}).Errorf("error getting user: %v", err)
		code, msg := customErr.ParseCustomError(err)
		errResp := customErr.NewErrorResponse(code, msg)
		c.JSON(http.StatusOK, errResp)
		return
	}

	c.JSON(http.StatusOK, res)
}
