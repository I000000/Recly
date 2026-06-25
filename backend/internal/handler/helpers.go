package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (string, bool) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return "", false
	}
	return userID, true
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
