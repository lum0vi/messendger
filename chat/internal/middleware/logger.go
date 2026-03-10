package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestID = "request_id"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New()
		c.Set(RequestID, requestId)
		logrus.WithFields(logrus.Fields{
			RequestID: requestId,
		}).Infof("request url %s", c.Request.RequestURI)

		c.Next()
	}
}
