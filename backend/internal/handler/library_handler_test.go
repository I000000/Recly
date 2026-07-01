package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupLibraryHandlerTest(t *testing.T) (*gin.Engine, *mocks.LibraryService, *LibraryHandler) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewLibraryService(t)
	handler := NewLibraryHandler(mockService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	return r, mockService, handler
}

func TestLibraryHandler_AddBook_Success(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.POST("/book/:id/like", handler.AddBook)

	mockService.On("AddBook", mock.Anything, "test-user-1", "book-123").Return(nil)

	req, _ := http.NewRequest("POST", "/book/book-123/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "book added to library", resp["message"])
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_AddBook_Error(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.POST("/book/:id/like", handler.AddBook)

	mockService.On("AddBook", mock.Anything, "test-user-1", "book-123").Return(assert.AnError)

	req, _ := http.NewRequest("POST", "/book/book-123/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], assert.AnError.Error())
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_RemoveBook_Success(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.DELETE("/book/:id/like", handler.RemoveBook)

	mockService.On("RemoveBook", mock.Anything, "test-user-1", "book-123").Return(nil)

	req, _ := http.NewRequest("DELETE", "/book/book-123/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "book removed from library", resp["message"])
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_GetBooks_Success(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.GET("/books", handler.GetBooks)

	expected := []domain.LikedBook{
		{UserID: "test-user-1", BookID: "book-1"},
		{UserID: "test-user-1", BookID: "book-2"},
	}
	mockService.On("GetBooks", mock.Anything, "test-user-1").Return(expected, nil)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	books := resp["books"].([]interface{})
	assert.Len(t, books, 2)
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_GetBooks_Error(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.GET("/books", handler.GetBooks)

	mockService.On("GetBooks", mock.Anything, "test-user-1").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_AddMovie_Success(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.POST("/movie/:id/like", handler.AddMovie)

	mockService.On("AddMovie", mock.Anything, "test-user-1", "movie-456").Return(nil)

	req, _ := http.NewRequest("POST", "/movie/movie-456/like", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "movie added to library", resp["message"])
	mockService.AssertExpectations(t)
}

func TestLibraryHandler_GetMovies_Success(t *testing.T) {
	r, mockService, handler := setupLibraryHandlerTest(t)
	r.GET("/movies", handler.GetMovies)

	expected := []domain.LikedMovie{
		{UserID: "test-user-1", MovieID: "movie-1"},
	}
	mockService.On("GetMovies", mock.Anything, "test-user-1").Return(expected, nil)

	req, _ := http.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	movies := resp["movies"].([]interface{})
	assert.Len(t, movies, 1)
	mockService.AssertExpectations(t)
}
