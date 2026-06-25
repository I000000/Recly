package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type SavedItemHandler struct {
	savedItemService interfaces.SavedItemService
}

func NewSavedItemHandler(savedItemService interfaces.SavedItemService) *SavedItemHandler {
	return &SavedItemHandler{savedItemService: savedItemService}
}

func (h *SavedItemHandler) Save(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req struct {
		ItemType string `json:"item_type" binding:"required,oneof=book movie"`
		ItemID   string `json:"item_id"   binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	item, err := h.savedItemService.SaveItem(c.Request.Context(), userID, req.ItemType, req.ItemID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"saved": item})
}

func (h *SavedItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.savedItemService.DeleteSavedItem(c.Request.Context(), id); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *SavedItemHandler) Get(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	items, err := h.savedItemService.GetSavedItems(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": items})
}
