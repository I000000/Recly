package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupViewedItemHandlerTest(t *testing.T) (*gin.Engine, *mocks.ViewedItemService, *ViewedItemHandler) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewViewedItemService(t)
	handler := NewViewedItemHandler(mockService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	return r, mockService, handler
}

func TestViewedItemHandler_RecordView_Success(t *testing.T) {
	r, mockService, handler := setupViewedItemHandlerTest(t)
	r.POST("/view", handler.RecordView)

	reqBody := recordViewRequest{
		ItemID:   "book-123",
		ItemType: "book",
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("RecordView", mock.Anything, "test-user-1", "book", "book-123").Return(nil)

	req, _ := http.NewRequest("POST", "/view", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
	mockService.AssertExpectations(t)
}

func TestViewedItemHandler_RecordView_BadRequest(t *testing.T) {
	r, mockService, handler := setupViewedItemHandlerTest(t)
	r.POST("/view", handler.RecordView)

	reqBody := map[string]string{
		"item_id":   "book-123",
		"item_type": "invalid",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/view", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "RecordView", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestViewedItemHandler_GetRecentViews_Success(t *testing.T) {
	r, mockService, handler := setupViewedItemHandlerTest(t)
	r.GET("/views", handler.GetRecentViews)

	now := time.Now()
	expected := []domain.ViewedItem{
		{ID: "view-1", UserID: "test-user-1", ItemType: "book", ItemID: "book-123", ViewedAt: now},
		{ID: "view-2", UserID: "test-user-1", ItemType: "movie", ItemID: "movie-456", ViewedAt: now.Add(-time.Hour)},
	}
	mockService.On("GetRecentViews", mock.Anything, "test-user-1", 20).Return(expected, nil)

	req, _ := http.NewRequest("GET", "/views", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	views := resp["views"].([]interface{})
	assert.Len(t, views, 2)
	mockService.AssertExpectations(t)
}

func TestViewedItemHandler_GetRecentViews_Error(t *testing.T) {
	r, mockService, handler := setupViewedItemHandlerTest(t)
	r.GET("/views", handler.GetRecentViews)

	mockService.On("GetRecentViews", mock.Anything, "test-user-1", 20).Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/views", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
