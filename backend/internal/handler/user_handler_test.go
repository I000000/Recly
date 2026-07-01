package handler

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserHandlerTest(t *testing.T) (*gin.Engine, *mocks.UserService, *UserHandler) {
	gin.SetMode(gin.TestMode)
	mockService := mocks.NewUserService(t)
	handler := NewUserHandler(mockService)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", "test-user-1")
		c.Next()
	})
	return r, mockService, handler
}

func TestUserHandler_Profile_Success(t *testing.T) {
	r, mockService, handler := setupUserHandlerTest(t)
	r.GET("/profile", handler.Profile)

	expectedUser := &domain.User{
		ID:                  "test-user-1",
		Email:               "test@example.com",
		Name:                "Test User",
		AvatarURL:           "http://example.com/avatar.png",
		OnboardingCompleted: true,
	}
	mockService.On("GetUserByID", mock.Anything, "test-user-1").Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "test-user-1", resp["id"])
	assert.Equal(t, "test@example.com", resp["email"])
	assert.Equal(t, "Test User", resp["name"])
	assert.Equal(t, "http://example.com/avatar.png", resp["avatar_url"])
	assert.Equal(t, true, resp["onboarding_completed"])
	mockService.AssertExpectations(t)
}

func TestUserHandler_Profile_NotFound(t *testing.T) {
	r, mockService, handler := setupUserHandlerTest(t)
	r.GET("/profile", handler.Profile)

	mockService.On("GetUserByID", mock.Anything, "test-user-1").Return(nil, domain.ErrNotFound)

	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "user not found")
	mockService.AssertExpectations(t)
}

func TestUserHandler_CompleteOnboarding_Success(t *testing.T) {
	r, mockService, handler := setupUserHandlerTest(t)
	r.POST("/onboarding/complete", handler.CompleteOnboarding)

	mockService.On("CompleteOnboarding", mock.Anything, "test-user-1").Return(nil)

	req, _ := http.NewRequest("POST", "/onboarding/complete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "ok", resp["status"])
	mockService.AssertExpectations(t)
}

func TestUserHandler_CompleteOnboarding_Error(t *testing.T) {
	r, mockService, handler := setupUserHandlerTest(t)
	r.POST("/onboarding/complete", handler.CompleteOnboarding)

	mockService.On("CompleteOnboarding", mock.Anything, "test-user-1").Return(assert.AnError)

	req, _ := http.NewRequest("POST", "/onboarding/complete", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUserHandler_UploadAvatar_Success(t *testing.T) {
	r, mockService, handler := setupUserHandlerTest(t)
	r.POST("/avatar", handler.UploadAvatar)

	mockService.On("UpdateAvatar", mock.Anything, "test-user-1", mock.Anything, mock.Anything).
		Return("http://example.com/new-avatar.png", nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("avatar", "avatar.png")
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/avatar", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "http://example.com/new-avatar.png", resp["avatar_url"])
	mockService.AssertExpectations(t)
}
