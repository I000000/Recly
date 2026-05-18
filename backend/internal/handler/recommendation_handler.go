package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recService *service.RecommendationService
}

func NewRecommendationHandler(recService *service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{recService: recService}
}

type RecommendRequest struct {
	SelectedIDs     []string           `json:"selected_ids" binding:"required,min=1"`
	Direction       string             `json:"direction" binding:"required,oneof=book_to_movie book_to_book movie_to_movie movie_to_book"`
	ModalityWeights map[string]float64 `json:"weights"`
}

func (h *RecommendationHandler) Request(c *gin.Context) {
	var req RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	taskID, err := h.recService.Request(c.Request.Context(), userID, req.SelectedIDs, req.Direction, req.ModalityWeights)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID, "status": "pending"})
}

func (h *RecommendationHandler) GetHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	history, err := h.recService.GetHistory(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": history})
}

type SaveRecommendationRequest struct {
	FromType string `json:"from_type" binding:"required,oneof=book movie"`
	FromID   string `json:"from_id" binding:"required"`
	ToType   string `json:"to_type" binding:"required,oneof=book movie"`
	ToID     string `json:"to_id" binding:"required"`
}

func (h *RecommendationHandler) Save(c *gin.Context) {
	var req SaveRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	rec, err := h.recService.SaveRecommendation(c.Request.Context(), userID, req.FromType, req.FromID, req.ToType, req.ToID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"saved": rec})
}

func (h *RecommendationHandler) DeleteSaved(c *gin.Context) {
	recID := c.Param("id")
	if err := h.recService.DeleteSavedRecommendation(c.Request.Context(), recID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *RecommendationHandler) GetSaved(c *gin.Context) {
	userID := c.GetString("user_id")
	saved, err := h.recService.GetSavedRecommendations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": saved})
}

func (h *RecommendationHandler) GetResult(c *gin.Context) {
	taskID := c.Param("taskId")
	result, err := h.recService.GetResult(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if result == nil {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	c.JSON(http.StatusOK, result)
}
