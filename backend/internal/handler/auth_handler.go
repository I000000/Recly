package handler

import (
	"errors"
	"net/http"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// Register handles user registration
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "user"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 409 {object} map[string]interface{} "error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid registration request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info("registering user", zap.String("email", req.Email))
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			log.Warn("registration failed: duplicate email", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
		log.Error("registration failed", zap.Error(err), zap.String("email", req.Email))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	log.Info("user registered successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "access_token, refresh_token"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid login request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Info("login attempt", zap.String("email", req.Email))
	access, refresh, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Warn("login failed", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	log.Info("login successful", zap.String("email", req.Email))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
