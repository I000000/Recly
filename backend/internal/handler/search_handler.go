package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/service"
	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{searchService: svc}
}

func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	itemType := c.DefaultQuery("type", "all") // book, movie, all
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}

func (h *SearchHandler) BatchGetItems(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids are required"})
		return
	}
	itemType := c.DefaultQuery("type", "movie")
	ids := strings.Split(idsStr, ",")

	items, err := h.searchService.GetItems(ids, itemType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
