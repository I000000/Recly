package handler

import (
	"net/http"

	"github.com/I000000/recly/internal/service/interfaces"
	"github.com/I000000/recly/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LibraryHandler struct {
	libService interfaces.LibraryService
}

func NewLibraryHandler(libService interfaces.LibraryService) *LibraryHandler {
	return &LibraryHandler{libService: libService}
}

// AddBook godoc
// @Summary Add book to library
// @Tags library
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /book/{id}/like [post]
func (h *LibraryHandler) AddBook(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	bookID := c.Param("id")
	log.Info("adding book to library", zap.String("user_id", userID), zap.String("book_id", bookID))
	if err := h.libService.AddBook(c.Request.Context(), userID, bookID); err != nil {
		log.Error("failed to add book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book added to library"})
}

// RemoveBook godoc
// @Summary Remove book from library
// @Tags library
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /book/{id}/like [delete]
func (h *LibraryHandler) RemoveBook(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	bookID := c.Param("id")
	log.Info("removing book from library", zap.String("user_id", userID), zap.String("book_id", bookID))
	if err := h.libService.RemoveBook(c.Request.Context(), userID, bookID); err != nil {
		log.Error("failed to remove book", zap.Error(err), zap.String("user_id", userID), zap.String("book_id", bookID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book removed from library"})
}

// GetBooks godoc
// @Summary Get liked books
// @Tags library
// @Produce json
// @Success 200 {object} map[string]interface{} "books"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/library/books [get]
func (h *LibraryHandler) GetBooks(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching liked books", zap.String("user_id", userID))
	books, err := h.libService.GetBooks(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to get books", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"books": books})
}

// AddMovie godoc
// @Summary Add movie to library
// @Tags library
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /movie/{id}/like [post]
func (h *LibraryHandler) AddMovie(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	movieID := c.Param("id")
	log.Info("adding movie to library", zap.String("user_id", userID), zap.String("movie_id", movieID))
	if err := h.libService.AddMovie(c.Request.Context(), userID, movieID); err != nil {
		log.Error("failed to add movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie added to library"})
}

// RemoveMovie godoc
// @Summary Remove movie from library
// @Tags library
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} map[string]interface{} "message"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /movie/{id}/like [delete]
func (h *LibraryHandler) RemoveMovie(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	movieID := c.Param("id")
	log.Info("removing movie from library", zap.String("user_id", userID), zap.String("movie_id", movieID))
	if err := h.libService.RemoveMovie(c.Request.Context(), userID, movieID); err != nil {
		log.Error("failed to remove movie", zap.Error(err), zap.String("user_id", userID), zap.String("movie_id", movieID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "movie removed from library"})
}

// GetMovies godoc
// @Summary Get liked movies
// @Tags library
// @Produce json
// @Success 200 {object} map[string]interface{} "movies"
// @Failure 401 {object} map[string]interface{} "unauthorized"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /user/library/movies [get]
func (h *LibraryHandler) GetMovies(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	log.Debug("fetching liked movies", zap.String("user_id", userID))
	movies, err := h.libService.GetMovies(c.Request.Context(), userID)
	if err != nil {
		log.Error("failed to get movies", zap.Error(err), zap.String("user_id", userID))
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"movies": movies})
}
