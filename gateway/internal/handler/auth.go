package handler

import (
	"gateway/internal/errors"
	"gateway/internal/middleware"
	"gateway/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Регистрирует нового пользователя. Принимает данные для регистрации и возвращает ответ с информацией о созданном пользователе.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.RegisterRequest  true  "Данные для регистрации нового пользователя"
// @Success      200  {object}  models.RegisterResponse  "Успешная регистрация"
// @Failure      400  {object}  errors.ErrorResponse  "Неверные данные для регистрации"
// @Failure      500  {object}  errors.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestId,
	}).Info("Register")

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestId,
		}).Errorf("Invalid register request: %v", err)
		errResp := errors.NewErrorResponse(http.StatusBadRequest, "Invalid register request")
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	res, err := h.service.Auth.Register(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestId,
		}).Error(err)
		code, err := errors.ParseCustomError(err)
		errResp := errors.NewErrorResponse(code, err)
		c.JSON(http.StatusBadRequest, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestId,
	}).Info("Register success")
	c.JSON(http.StatusOK, res)
}

// Login godoc
// @Security BearerAuth
// @Summary      Авторизация пользователя
// @Description  Выполняет авторизацию пользователя и генерирует JWT токен. Принимает данные для входа и возвращает токен при успешной авторизации.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  models.LoginRequest  true  "Данные для авторизации пользователей"
// @Success      200  {object}  models.LoginResponse  "Успешная авторизация, возвращает JWT токен"
// @Failure      400  {object}  errors.ErrorResponse  "Неверные данные для авторизации"
// @Failure      401  {object}  errors.ErrorResponse  "Не авторизован, неверные учетные данные"
// @Failure      500  {object}  errors.ErrorResponse  "Внутренняя ошибка сервера"
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	requestId, ok := c.Get(middleware.RequestIDKey)
	if !ok {
		requestId = "unknown"
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestId,
	}).Info("Login")
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestId,
		}).Errorf("Invalid login request: %v", err)
		errResp := errors.NewErrorResponse(http.StatusBadRequest, err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestId,
	}).Infof("request: %+v", req)
	res, err := h.service.Auth.Login(&req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			middleware.RequestIDKey: requestId,
		}).Error(err)
		code, err := errors.ParseCustomError(err)
		errResp := errors.NewErrorResponse(code, err)
		c.AbortWithStatusJSON(code, errResp)
		return
	}
	logrus.WithFields(logrus.Fields{
		middleware.RequestIDKey: requestId,
	}).Info("Login success")
	c.JSON(http.StatusOK, res)
}
