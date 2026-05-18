package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerWithoutSpam() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/announce") ||
			strings.HasPrefix(path, "/control") ||
			strings.HasPrefix(path, "/ann") {
			c.AbortWithStatus(403)
			return
		}

		start := time.Now()
		c.Next()
		latency := time.Since(start)
		gin.LoggerWithWriter(gin.DefaultWriter, "/health")(c)
		_ = latency
	}
}
