package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Profile godoc
// @Summary Get user profile
// @Tags user
// @Produce json
// @Success 200 {object} domain.User "user profile"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 404 {object} map[string]interface{} "user not found"
// @Router /user/profile [get]
func (h *UserHandler) Profile(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching user profile", zap.String("user_id", userID))
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Warn("user profile not found", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusNotFound, "user not found")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":                   user.ID,
		"email":                user.Email,
		"name":                 user.Name,
		"avatar_url":           user.AvatarURL,
		"onboarding_completed": user.OnboardingCompleted,
	})
}

// CompleteOnboarding godoc
// @Summary Complete onboarding
// @Tags user
// @Produce json
// @Success 200 {object} map[string]interface{} "status"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/onboarding/complete [post]
func (h *UserHandler) CompleteOnboarding(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Info("completing onboarding", zap.String("user_id", userID))
	if err := h.userService.CompleteOnboarding(c.Request.Context(), userID); err != nil {
		log.Error("failed to complete onboarding", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UploadAvatar godoc
// @Summary Upload avatar
// @Tags user
// @Accept mpfd
// @Produce json
// @Param avatar formData file true "Avatar image (jpg, png, webp)"
// @Success 200 {object} map[string]interface{} "avatar_url"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		log.Warn("no avatar file provided", zap.Error(err))
		respondWithError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	const maxSize = 5 * 1024 * 1024
	if header.Size > maxSize {
		log.Warn("avatar file too large", zap.Int64("size", header.Size), zap.String("user_id", userID))
		respondWithError(c, http.StatusBadRequest, "file too large (max 5MB)")
		return
	}
	log.Info("uploading avatar", zap.String("user_id", userID), zap.String("filename", header.Filename), zap.Int64("size", header.Size))
	avatarURL, err := h.userService.UpdateAvatar(c.Request.Context(), userID, file, header)
	if err != nil {
		log.Error("failed to upload avatar", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("avatar uploaded successfully", zap.String("user_id", userID), zap.String("avatar_url", avatarURL))
	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
