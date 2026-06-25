package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
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

func (h *UserHandler) CompleteOnboarding(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	if err := h.userService.CompleteOnboarding(c.Request.Context(), userID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	const maxSize = 5 * 1024 * 1024
	if header.Size > maxSize {
		respondWithError(c, http.StatusBadRequest, "file too large (max 5MB)")
		return
	}

	avatarURL, err := h.userService.UpdateAvatar(c.Request.Context(), userID, file, header)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"avatar_url": avatarURL})
}
