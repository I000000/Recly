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

// Save saves an item to saved list
// @Summary Save item
// @Tags saved
// @Accept json
// @Produce json
// @Param request body object true "Item data (item_type, item_id)"
// @Success 201 {object} map[string]interface{} "saved"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [post]
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

// Delete deletes a saved item
// @Summary Delete saved item
// @Tags saved
// @Produce json
// @Param id path string true "Saved item ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items/{id} [delete]
func (h *SavedItemHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.savedItemService.DeleteSavedItem(c.Request.Context(), id); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Get returns user's saved items
// @Summary Get saved items
// @Tags saved
// @Produce json
// @Success 200 {object} map[string]interface{} "saved"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [get]
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
