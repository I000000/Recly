package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SpamFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/announce") ||
			strings.HasPrefix(path, "/control") ||
			strings.HasPrefix(path, "/ann") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
