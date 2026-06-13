package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service"
	"github.com/gin-gonic/gin"
)

type ViewedItemHandler struct {
	viewedItemService *service.ViewedItemService
}

func NewViewedItemHandler(svc *service.ViewedItemService) *ViewedItemHandler {
	return &ViewedItemHandler{viewedItemService: svc}
}

type recordViewRequest struct {
	ItemID   string `json:"item_id" binding:"required"`
	ItemType string `json:"item_type" binding:"required,oneof=book movie"`
}

func (h *ViewedItemHandler) RecordView(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req recordViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.viewedItemService.RecordView(c.Request.Context(), userID, req.ItemType, req.ItemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ViewedItemHandler) GetRecentViews(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit := 20 // можно вынести в параметр запроса
	views, err := h.viewedItemService.GetRecentViews(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"views": views})
}
