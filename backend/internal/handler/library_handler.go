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

// AddBook adds a book to user's library
// @Summary Add book to library
// @Tags library
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /book/{id}/like [post]
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

// RemoveBook removes a book from user's library
// @Summary Remove book from library
// @Tags library
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /book/{id}/like [delete]
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

// GetBooks returns user's liked books
// @Summary Get liked books
// @Tags library
// @Produce json
// @Success 200 {object} map[string]interface{} "books"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/library/books [get]
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

// AddMovie adds a movie to user's library
// @Summary Add movie to library
// @Tags library
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /movie/{id}/like [post]
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

// RemoveMovie removes a movie from user's library
// @Summary Remove movie from library
// @Tags library
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /movie/{id}/like [delete]
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

// GetMovies returns user's liked movies
// @Summary Get liked movies
// @Tags library
// @Produce json
// @Success 200 {object} map[string]interface{} "movies"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/library/movies [get]
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
