package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestIDKey = "requestID"

func RequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New()
		c.Set(RequestIDKey, requestId)
		logrus.WithFields(logrus.Fields{
			RequestIDKey: requestId,
		}).Infof("request url %s", c.Request.RequestURI)

		c.Next()
	}
}
