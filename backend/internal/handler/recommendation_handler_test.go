package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/redis"
	"github.com/I000000/recly/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRecommendationHandlerTest(t *testing.T) (*gin.Engine, *mocks.RecommendationService, *RecommendationHandler) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewRecommendationService(t)
	handler := NewRecommendationHandler(mockService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	return r, mockService, handler
}

func TestRecommendationHandler_Request_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.POST("/recommend", handler.Request)

	reqBody := RecommendRequest{
		SelectedIDs:     []string{"book_1", "movie_1"},
		ExcludeIDs:      []string{},
		Direction:       "book_to_movie",
		ModalityWeights: map[string]float64{"text": 0.4, "genre": 0.3, "image": 0.3},
		Contextual:      false,
	}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("Request", mock.Anything, "test-user-1", reqBody.SelectedIDs, reqBody.ModalityWeights, reqBody.ExcludeIDs, reqBody.Direction, reqBody.Contextual).
		Return("task-123", nil)

	req, _ := http.NewRequest("POST", "/recommend", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "task-123", resp["task_id"])
	assert.Equal(t, "pending", resp["status"])
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_Request_BadRequest(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.POST("/recommend", handler.Request)

	reqBody := map[string]interface{}{
		"selected_ids": []string{"book_1"},
		"direction":    "invalid_direction",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/recommend", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "Request", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestRecommendationHandler_GetHistory_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.GET("/history", handler.GetHistory)

	expectedHistory := []domain.RecommendationHistory{
		{
			ID:        "hist-1",
			UserID:    "test-user-1",
			TaskID:    "task-1",
			Direction: "book_to_movie",
			Result:    `["movie_1"]`,
		},
	}
	mockService.On("GetHistory", mock.Anything, "test-user-1").Return(expectedHistory, nil)

	req, _ := http.NewRequest("GET", "/history", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	history := resp["history"].([]interface{})
	assert.Len(t, history, 1)
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_GetHistory_Error(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.GET("/history", handler.GetHistory)

	mockService.On("GetHistory", mock.Anything, "test-user-1").Return(nil, assert.AnError)

	req, _ := http.NewRequest("GET", "/history", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_Save_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.POST("/saved", handler.Save)

	reqBody := SaveRecommendationRequest{
		FromType: "book",
		FromID:   "book_1",
		ToType:   "movie",
		ToID:     "movie_1",
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedSaved := &domain.SavedRecommendation{
		ID:       "saved-1",
		UserID:   "test-user-1",
		FromType: "book",
		FromID:   "book_1",
		ToType:   "movie",
		ToID:     "movie_1",
	}
	mockService.On("SaveRecommendation", mock.Anything, "test-user-1", reqBody.FromType, reqBody.FromID, reqBody.ToType, reqBody.ToID).
		Return(expectedSaved, nil)

	req, _ := http.NewRequest("POST", "/saved", bytes.NewReader(jsonBody))
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

func TestRecommendationHandler_DeleteSaved_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.DELETE("/saved/:id", handler.DeleteSaved)

	mockService.On("DeleteSavedRecommendation", mock.Anything, "saved-1").Return(nil)

	req, _ := http.NewRequest("DELETE", "/saved/saved-1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "deleted", resp["message"])
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_GetSaved_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.GET("/saved", handler.GetSaved)

	expectedSaved := []domain.SavedRecommendation{
		{ID: "saved-1", UserID: "test-user-1", FromType: "book", FromID: "book_1", ToType: "movie", ToID: "movie_1"},
	}
	mockService.On("GetSavedRecommendations", mock.Anything, "test-user-1").Return(expectedSaved, nil)

	req, _ := http.NewRequest("GET", "/saved", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	saved := resp["saved"].([]interface{})
	assert.Len(t, saved, 1)
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_GetResult_Success(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.GET("/result/:taskId", handler.GetResult)

	expectedResult := &redis.RecommendationResult{
		Status: "done",
		Movies: []string{"movie_1", "movie_2"},
	}
	mockService.On("GetResult", mock.Anything, "task-123").Return(expectedResult, nil)

	req, _ := http.NewRequest("GET", "/result/task-123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "done", resp["status"])
	movies := resp["movies"].([]interface{})
	assert.Len(t, movies, 2)
	mockService.AssertExpectations(t)
}

func TestRecommendationHandler_GetResult_Pending(t *testing.T) {
	r, mockService, handler := setupRecommendationHandlerTest(t)
	r.GET("/result/:taskId", handler.GetResult)

	mockService.On("GetResult", mock.Anything, "task-pending").Return(nil, nil)

	req, _ := http.NewRequest("GET", "/result/task-pending", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "pending", resp["status"])
	mockService.AssertExpectations(t)
}
