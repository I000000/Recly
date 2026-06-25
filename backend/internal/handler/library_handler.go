package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	libService interfaces.LibraryService
}

func NewLibraryHandler(libService interfaces.LibraryService) *LibraryHandler {
	return &LibraryHandler{libService: libService}
}

func (h *LibraryHandler) AddBook(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	bookID := c.Param("id")
	if err := h.libService.AddBook(c.Request.Context(), userID, bookID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book added to library"})
}

func (h *LibraryHandler) RemoveBook(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	bookID := c.Param("id")
	if err := h.libService.RemoveBook(c.Request.Context(), userID, bookID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book removed from library"})
}

func (h *LibraryHandler) GetBooks(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	books, err := h.libService.GetBooks(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

func (h *LibraryHandler) AddMovie(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	movieID := c.Param("id")
	if err := h.libService.AddMovie(c.Request.Context(), userID, movieID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie added to library"})
}

func (h *LibraryHandler) RemoveMovie(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	movieID := c.Param("id")
	if err := h.libService.RemoveMovie(c.Request.Context(), userID, movieID); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie removed from library"})
}

func (h *LibraryHandler) GetMovies(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	movies, err := h.libService.GetMovies(c.Request.Context(), userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
