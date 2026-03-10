package middleware

import (
	"net/http"

	"gateway/internal/jwt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const UserIDKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId, ok := c.Get(RequestIDKey)
		if !ok {
			requestId = "unknown"
		}

		// Пробуем получить токен из заголовка или из query параметра (для WebSocket)
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		// Убираем префикс "Bearer " если есть
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if token == "" {
			logrus.WithFields(logrus.Fields{
				RequestIDKey: requestId,
			}).Info("empty token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userId, err := jwt.ParseJWT(token)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				RequestIDKey: requestId,
			}).Infof("invalid token: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		logrus.WithFields(logrus.Fields{
			UserIDKey: userId,
		}).Infof("authorized user id: %v", userId)
		c.Set(UserIDKey, userId)
		c.Next()
	}
}
