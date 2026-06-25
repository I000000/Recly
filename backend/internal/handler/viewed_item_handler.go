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
