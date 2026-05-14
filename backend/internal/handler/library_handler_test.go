package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I000000/recly/internal/handler"
	"github.com/I000000/recly/internal/service"
	"github.com/I000000/recly/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestLibraryHandler_AddBook(t *testing.T) {
	repo := &mocks.LibraryRepository{}
	repo.On("AddLikedBook", mock.Anything, "user1", "book1").Return(nil)
	svc := service.NewLibraryService(repo)
	h := handler.NewLibraryHandler(svc)

	r := setupTestRouter()
	r.POST("/book/:id/like", func(c *gin.Context) { c.Set("user_id", "user1"); h.AddBook(c) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/book/book1/like", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "book added to library", resp["message"])
	repo.AssertExpectations(t)
}
