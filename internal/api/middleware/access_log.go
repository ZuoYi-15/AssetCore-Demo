package middleware

import (
	"time"

	"asset-core/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func AccessLog(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http request",
			logger.String("method", c.Request.Method),
			logger.String("path", c.Request.URL.Path),
			logger.Int("status", c.Writer.Status()),
			logger.String("trace_id", c.GetString("trace_id")),
			logger.Int("latency_ms", int(time.Since(start).Milliseconds())),
		)
	}
}
