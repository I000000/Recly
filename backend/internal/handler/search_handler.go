package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService interfaces.SearchService
}

func NewSearchHandler(searchService interfaces.SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// Search performs full-text search
// @Summary Search items
// @Tags search
// @Produce json
// @Param q query string false "Search query"
// @Param type query string false "Item type (book/movie/all)" default(all)
// @Param genre query string false "Genre filter"
// @Param sort query string false "Sort field: ratings_count:desc, vote_count:desc"
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} map[string]interface{} "results"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	itemType := c.DefaultQuery("type", "all")
	genre := c.Query("genre")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if query == "" && genre == "" && (itemType == "" || itemType == "all") {
		c.JSON(http.StatusOK, gin.H{"results": []domain.ItemDetail{}})
		return
	}

	results, err := h.searchService.SearchWithFilters(query, itemType, genre, sort, limit, offset)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

// BatchGetItems returns metadata for multiple items
// @Summary Get multiple items by IDs
// @Tags search
// @Produce json
// @Param ids query string true "Comma-separated IDs"
// @Param type query string false "Item type (book/movie/all)" default(all)
// @Success 200 {object} map[string]interface{} "items"
// @Failure 400 {object} map[string]interface{} "error"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /items/batch [get]
func (h *SearchHandler) BatchGetItems(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		respondWithError(c, http.StatusBadRequest, "ids are required")
		return
	}
	itemType := c.DefaultQuery("type", "movie")
	ids := strings.Split(idsStr, ",")

	items, err := h.searchService.GetItems(ids, itemType)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// GetGenres returns list of genres
// @Summary Get genres
// @Tags search
// @Produce json
// @Param type query string false "Item type (book/movie/all)" default(all)
// @Success 200 {object} map[string]interface{} "genres"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /genres [get]
func (h *SearchHandler) GetGenres(c *gin.Context) {
	itemType := c.DefaultQuery("type", "all")
	genres, err := h.searchService.GetGenres(itemType)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"genres": genres})
}
