package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	recService interfaces.RecommendationService
}

func NewRecommendationHandler(recService interfaces.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{recService: recService}
}

type RecommendRequest struct {
	SelectedIDs     []string           `json:"selected_ids" binding:"required"`
	ExcludeIDs      []string           `json:"exclude_ids"`
	Direction       string             `json:"direction" binding:"required,oneof=book_to_movie book_to_book movie_to_movie movie_to_book"`
	ModalityWeights map[string]float64 `json:"weights"`
	Contextual      bool               `json:"contextual"`
}

// Request creates a recommendation task
// @Summary Request recommendations
// @Tags recommendations
// @Accept json
// @Produce json
// @Param request body RecommendRequest true "Recommendation request"
// @Success 202 {object} map[string]interface{} "task_id, status"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /recommend [post]
func (h *RecommendationHandler) Request(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	taskID, err := h.recService.Request(c.Request.Context(), userID, req.SelectedIDs, req.ModalityWeights, req.ExcludeIDs, req.Direction, req.Contextual)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID, "status": "pending"})
}

// GetResult returns recommendation result by task ID
// @Summary Get recommendation result
// @Tags recommendations
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} redis.RecommendationResult
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /result/{taskId} [get]
func (h *RecommendationHandler) GetResult(c *gin.Context) {
	taskID := c.Param("taskId")
	result, err := h.recService.GetResult(c.Request.Context(), taskID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "internal error")
		return
	}
	if result == nil {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetHistory returns user's recommendation history
// @Summary Get recommendation history
// @Tags recommendations
// @Produce json
// @Success 200 {object} map[string]interface{} "history"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/recommendations/history [get]
func (h *RecommendationHandler) GetHistory(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	history, err := h.recService.GetHistory(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
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

// Save saves a recommendation to user's saved list
// @Summary Save recommendation
// @Tags recommendations
// @Accept json
// @Produce json
// @Param request body SaveRecommendationRequest true "Save data"
// @Success 201 {object} map[string]interface{} "saved"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [post]
func (h *RecommendationHandler) Save(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req SaveRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	rec, err := h.recService.SaveRecommendation(c.Request.Context(), userID, req.FromType, req.FromID, req.ToType, req.ToID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"saved": rec})
}

// DeleteSaved deletes a saved recommendation
// @Summary Delete saved recommendation
// @Tags recommendations
// @Produce json
// @Param id path string true "Saved item ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items/{id} [delete]
func (h *RecommendationHandler) DeleteSaved(c *gin.Context) {
	recID := c.Param("id")
	if err := h.recService.DeleteSavedRecommendation(c.Request.Context(), recID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetSaved returns user's saved recommendations
// @Summary Get saved recommendations
// @Tags recommendations
// @Produce json
// @Success 200 {object} map[string]interface{} "saved"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [get]
func (h *RecommendationHandler) GetSaved(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	saved, err := h.recService.GetSavedRecommendations(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": saved})
}
