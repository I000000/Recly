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

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockAuthService := mocks.NewAuthService(t)

	mockAuthService.On("Register", mock.Anything, "test@example.com", "password123", "Test User").
		Return(&domain.User{
			ID:    "user-1",
			Email: "test@example.com",
			Name:  "Test User",
		}, nil)

	handler := NewAuthHandler(mockAuthService)

	r := setupTestRouter()
	r.POST("/register", handler.Register)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "user")
	user := resp["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", user["email"])
	assert.Equal(t, "Test User", user["name"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {
	mockAuthService := mocks.NewAuthService(t)
	handler := NewAuthHandler(mockAuthService)

	r := setupTestRouter()
	r.POST("/register", handler.Register)

	body := map[string]string{
		"email":    "invalid",
		"password": "password123",
		"name":     "Test",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockAuthService.AssertNotCalled(t, "Register", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestAuthHandler_Register_DuplicateEmail(t *testing.T) {
	mockAuthService := mocks.NewAuthService(t)
	mockAuthService.On("Register", mock.Anything, "test@example.com", "password123", "Test User").
		Return(nil, domain.ErrDuplicateEmail)

	handler := NewAuthHandler(mockAuthService)

	r := setupTestRouter()
	r.POST("/register", handler.Register)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
		"name":     "Test User",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "email already exists")

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockAuthService := mocks.NewAuthService(t)
	mockAuthService.On("Login", mock.Anything, "test@example.com", "password123").
		Return("access-token", "refresh-token", nil)

	handler := NewAuthHandler(mockAuthService)

	r := setupTestRouter()
	r.POST("/login", handler.Login)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp, "access_token")
	assert.Contains(t, resp, "refresh_token")
	assert.Equal(t, "access-token", resp["access_token"])

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockAuthService := mocks.NewAuthService(t)
	mockAuthService.On("Login", mock.Anything, "test@example.com", "wrongpass").
		Return("", "", assert.AnError)

	handler := NewAuthHandler(mockAuthService)

	r := setupTestRouter()
	r.POST("/login", handler.Login)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpass",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "invalid credentials")

	mockAuthService.AssertExpectations(t)
}
