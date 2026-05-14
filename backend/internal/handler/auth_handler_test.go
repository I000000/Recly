package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/I000000/recly/internal/domain"
	"github.com/I000000/recly/internal/handler"
)

// Mock для интерфейса AuthServiceInterface
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Register(ctx context.Context, email, password, name string) (*domain.User, error) {
	args := m.Called(ctx, email, password, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.Error(2)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	mockSvc := new(mockAuthService)
	mockSvc.On("Register", mock.Anything, "test@example.com", "password123", "Test User").
		Return(&domain.User{ID: "user1", Email: "test@example.com", Name: "Test User"}, nil)

	h := handler.NewAuthHandler(mockSvc)
	r := gin.New()
	r.POST("/register", h.Register)

	body := `{"email":"test@example.com","password":"password123","name":"Test User"}`
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	user := resp["user"].(map[string]interface{})
	assert.Equal(t, "test@example.com", user["email"])
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := new(mockAuthService)
	mockSvc.On("Login", mock.Anything, "test@example.com", "password123").
		Return("access-token", "refresh-token", nil)

	h := handler.NewAuthHandler(mockSvc)
	r := gin.New()
	r.POST("/login", h.Login)

	body := `{"email":"test@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "access-token", resp["access_token"])
	mockSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	mockSvc := new(mockAuthService)
	mockSvc.On("Login", mock.Anything, "test@example.com", "wrong").Return("", "", errors.New("invalid credentials"))

	h := handler.NewAuthHandler(mockSvc)
	r := gin.New()
	r.POST("/login", h.Login)

	body := `{"email":"test@example.com","password":"wrong"}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockSvc.AssertExpectations(t)
}
