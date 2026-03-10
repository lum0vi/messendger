package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestID = "RequestID"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Set(RequestID, requestId)
		logrus.WithFields(logrus.Fields{
			RequestID: requestId,
		}).Info("request started")
		c.Next()
	}
}
