package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service"
	"github.com/gin-gonic/gin"
)

type SavedItemHandler struct {
	savedItemService *service.SavedItemService
}

func NewSavedItemHandler(svc *service.SavedItemService) *SavedItemHandler {
	return &SavedItemHandler{savedItemService: svc}
}

func (h *SavedItemHandler) Save(c *gin.Context) {
	var req struct {
		ItemType string `json:"item_type" binding:"required,oneof=book movie"`
		ItemID   string `json:"item_id"   binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	item, err := h.savedItemService.SaveItem(c.Request.Context(), userID, req.ItemType, req.ItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"saved": item})
}

func (h *SavedItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.savedItemService.DeleteSavedItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *SavedItemHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	items, err := h.savedItemService.GetSavedItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": items})
}
