package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type ViewedItemHandler struct {
	viewedItemService interfaces.ViewedItemService
}

func NewViewedItemHandler(viewedItemService interfaces.ViewedItemService) *ViewedItemHandler {
	return &ViewedItemHandler{viewedItemService: viewedItemService}
}

type recordViewRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	ItemType string `json:"item_type" binding:"required,oneof=book movie"`
}

// RecordView records a user's view of an item
// @Summary Record view
// @Tags views
// @Accept json
// @Produce json
// @Param request body recordViewRequest true "View data"
// @Success 200 {object} map[string]interface{} "status"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/view [post]
func (h *ViewedItemHandler) RecordView(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req recordViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.viewedItemService.RecordView(c.Request.Context(), userID, req.ItemType, req.ItemID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetRecentViews returns user's recent views
// @Summary Get recent views
// @Tags views
// @Produce json
// @Success 200 {object} map[string]interface{} "views"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/views [get]
func (h *ViewedItemHandler) GetRecentViews(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	limit := 20
	views, err := h.viewedItemService.GetRecentViews(c.Request.Context(), userID, limit)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"views": views})
}
