package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetString("user_id")
	// позже можно получить email, name из БД
	c.JSON(http.StatusOK, gin.H{"user_id": userID, "name": "", "email": ""})
}
