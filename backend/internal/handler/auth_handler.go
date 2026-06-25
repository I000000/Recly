package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService interfaces.AuthService
}

func NewAuthHandler(authService interfaces.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
