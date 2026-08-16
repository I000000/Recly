package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// Request godoc
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
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req RecommendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid recommendation request", zap.Error(err))
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Info("requesting recommendations",
		zap.String("user_id", userID),
		zap.String("direction", req.Direction),
		zap.Int("selected_count", len(req.SelectedIDs)),
		zap.Bool("contextual", req.Contextual),
	)
	taskID, err := h.recService.Request(c.Request.Context(), userID, req.SelectedIDs, req.ModalityWeights, req.ExcludeIDs, req.Direction, req.Contextual)
	if err != nil {
		log.Error("failed to request recommendations", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("recommendation task created", zap.String("task_id", taskID), zap.String("user_id", userID))
	c.JSON(http.StatusAccepted, gin.H{"task_id": taskID, "status": "pending"})
}

// GetHistory godoc
// @Summary Get recommendation history
// @Tags recommendations
// @Produce json
// @Success 200 {object} map[string]interface{} "history"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/recommendations/history [get]
func (h *RecommendationHandler) GetHistory(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching recommendation history", zap.String("user_id", userID))
	history, err := h.recService.GetHistory(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to get history", zap.Error(err), zap.String("user_id", userID))
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

// Save godoc
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
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req SaveRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid save request", zap.Error(err))
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Info("saving recommendation",
		zap.String("user_id", userID),
		zap.String("from_type", req.FromType),
		zap.String("from_id", req.FromID),
		zap.String("to_type", req.ToType),
		zap.String("to_id", req.ToID),
	)
	rec, err := h.recService.SaveRecommendation(c.Request.Context(), userID, req.FromType, req.FromID, req.ToType, req.ToID)
	if err != nil {
		log.Error("failed to save recommendation", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("recommendation saved", zap.String("user_id", userID), zap.String("saved_id", rec.ID))
	c.JSON(http.StatusCreated, gin.H{"saved": rec})
}

// DeleteSaved godoc
// @Summary Delete saved recommendation
// @Tags recommendations
// @Produce json
// @Param id path string true "Saved item ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items/{id} [delete]
func (h *RecommendationHandler) DeleteSaved(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	recID := c.Param("id")
	log.Info("deleting saved recommendation", zap.String("id", recID))
	if err := h.recService.DeleteSavedRecommendation(c.Request.Context(), recID); err != nil {
		log.Error("failed to delete saved recommendation", zap.Error(err), zap.String("id", recID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// GetSaved godoc
// @Summary Get saved recommendations
// @Tags recommendations
// @Produce json
// @Success 200 {object} map[string]interface{} "saved"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [get]
func (h *RecommendationHandler) GetSaved(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching saved recommendations", zap.String("user_id", userID))
	saved, err := h.recService.GetSavedRecommendations(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to get saved recommendations", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": saved})
}

// GetResult godoc
// @Summary Get recommendation result
// @Tags recommendations
// @Produce json
// @Param taskId path string true "Task ID"
// @Success 200 {object} redis.RecommendationResult
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /result/{taskId} [get]
func (h *RecommendationHandler) GetResult(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	taskID := c.Param("taskId")
	log.Debug("fetching recommendation result", zap.String("task_id", taskID))
	result, err := h.recService.GetResult(c.Request.Context(), taskID)
	if err != nil {
		log.Error("failed to get result", zap.Error(err), zap.String("task_id", taskID))
		respondWithError(c, http.StatusInternalServerError, "internal error")
		return
	}
	if result == nil {
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	c.JSON(http.StatusOK, result)
}
