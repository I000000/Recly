package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SavedItemHandler struct {
	savedItemService interfaces.SavedItemService
}

func NewSavedItemHandler(savedItemService interfaces.SavedItemService) *SavedItemHandler {
	return &SavedItemHandler{savedItemService: savedItemService}
}

// Save godoc
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
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	var req struct {
		ItemType string `json:"item_type" binding:"required,oneof=book movie"`
		ItemID   string `json:"item_id"   binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("invalid saved item request", zap.Error(err))
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	log.Info("saving item", zap.String("user_id", userID), zap.String("item_type", req.ItemType), zap.String("item_id", req.ItemID))
	item, err := h.savedItemService.SaveItem(c.Request.Context(), userID, req.ItemType, req.ItemID)
	if err != nil {
		log.Error("failed to save item", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	log.Info("item saved", zap.String("user_id", userID), zap.String("saved_id", item.ID))
	c.JSON(http.StatusCreated, gin.H{"saved": item})
}

// Delete godoc
// @Summary Delete saved item
// @Tags saved
// @Produce json
// @Param id path string true "Saved item ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items/{id} [delete]
func (h *SavedItemHandler) Delete(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	id := c.Param("id")
	log.Info("deleting saved item", zap.String("id", id))
	if err := h.savedItemService.DeleteSavedItem(c.Request.Context(), id); err != nil {
		log.Error("failed to delete saved item", zap.Error(err), zap.String("id", id))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Get godoc
// @Summary Get saved items
// @Tags saved
// @Produce json
// @Success 200 {object} map[string]interface{} "saved"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/saved-items [get]
func (h *SavedItemHandler) Get(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching saved items", zap.String("user_id", userID))
	items, err := h.savedItemService.GetSavedItems(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to get saved items", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"saved": items})
}
