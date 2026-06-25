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

func (h *SearchHandler) GetGenres(c *gin.Context) {
	itemType := c.DefaultQuery("type", "all")
	genres, err := h.searchService.GetGenres(itemType)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"genres": genres})
}
