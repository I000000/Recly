package handler

import (
	"bytes"
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

func setupSavedItemHandlerTest(t *testing.T) (*gin.Engine, *mocks.SavedItemService, *SavedItemHandler) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewSavedItemService(t)
	handler := NewSavedItemHandler(mockService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	return r, mockService, handler
}

func TestSavedItemHandler_Save_Success(t *testing.T) {
	r, mockService, handler := setupSavedItemHandlerTest(t)
	r.POST("/saved-items", handler.Save)

	reqBody := struct {
		ItemType string `json:"item_type"`
		ItemID   string `json:"item_id"`
	}{
		ItemType: "book",
		ItemID:   "book-123",
	}
	jsonBody, _ := json.Marshal(reqBody)

	expected := &domain.SavedItem{
		ID:       "saved-1",
		UserID:   "test-user-1",
		ItemType: "book",
		ItemID:   "book-123",
	}
	mockService.On("SaveItem", mock.Anything, "test-user-1", "book", "book-123").Return(expected, nil)

	req, _ := http.NewRequest("POST", "/saved-items", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	saved := resp["saved"].(map[string]interface{})
	assert.Equal(t, "saved-1", saved["id"])
	mockService.AssertExpectations(t)
}

func TestSavedItemHandler_Delete_Success(t *testing.T) {
	r, mockService, handler := setupSavedItemHandlerTest(t)
	r.DELETE("/saved-items/:id", handler.Delete)

	mockService.On("DeleteSavedItem", mock.Anything, "saved-1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/saved-items/saved-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "deleted", resp["message"])
	mockService.AssertExpectations(t)
}

func TestSavedItemHandler_Get_Success(t *testing.T) {
	r, mockService, handler := setupSavedItemHandlerTest(t)
	r.GET("/saved-items", handler.Get)

	expected := []domain.SavedItem{
		{ID: "saved-1", UserID: "test-user-1", ItemType: "book", ItemID: "book-123"},
	}
	mockService.On("GetSavedItems", mock.Anything, "test-user-1").Return(expected, nil)

	req, _ := http.NewRequest("GET", "/saved-items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["saved"].([]interface{})
	assert.Len(t, items, 1)
	mockService.AssertExpectations(t)
}

func TestSavedItemHandler_Get_Error(t *testing.T) {
	r, mockService, handler := setupSavedItemHandlerTest(t)
	r.GET("/saved-items", handler.Get)

	mockService.On("GetSavedItems", mock.Anything, "test-user-1").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/saved-items", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
