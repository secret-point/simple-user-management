package middlewares

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// RequestLogger returns a gin.HandlerFunc (middleware) that logs requests using Logrus.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// End timer
		end := time.Now()
		latency := end.Sub(start)

		// Log details
		log.WithFields(log.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"ip":         c.ClientIP(),
			"latency":    latency,
			"user-agent": c.Request.UserAgent(),
		}).Info("Request received")
	}
}
