package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sundayyogurt/user_service/internal/domain"
	"github.com/sundayyogurt/user_service/internal/dto"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(input dto.UserSignup) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserService) Login(email, password string) (*domain.User, error) {
	args := m.Called(email, password)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserService) ForgotPassword(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserService) SetPassword(token, password string) error {
	args := m.Called(token, password)
	return args.Error(0)
}

func (m *MockUserService) CreateProfile(profile dto.UserProfile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockUserService) GetProfile(userID int) (*domain.User, error) {
	args := m.Called(userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserService) Authenticate(ctx *fiber.Ctx) (*domain.User, error) {
	args := m.Called(ctx)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func setupTest(t *testing.T) (*fiber.App, *MockUserService, *UserHandler) {
	app := fiber.New()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	return app, mockService, handler
}

func TestUserHandler_Register(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/register", handler.Register)

	body := dto.UserSignup{Email: "test@test.com", Password: "password", Phone: "123456789"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("Register", mock.Anything).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestRegister_InvalidBody(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/register", handler.Register)

	body := dto.UserSignup{Email: "", Password: "password", Phone: "123456789"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("Register", mock.Anything).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRegister_ServiceError(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/register", handler.Register)

	body := dto.UserSignup{Email: "tset@test.com", Password: "password", Phone: "123456789"}
	BodyJSON, _ := json.Marshal(body)
	mockService.On("Register", mock.Anything).Return(fiber.NewError(http.StatusInternalServerError, "Internal Server Error"))
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(BodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/login", handler.Login)

	body := dto.UserLogin{Email: "test@test.com", Password: "password"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("Login", body.Email, body.Password).Return(&domain.User{ID: 1, Email: body.Email}, nil)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLogin_UnAuthorized(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/login", handler.Login)

	body := dto.UserLogin{Email: "test@test.com", Password: "wrong_password"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("Login", body.Email, body.Password).Return(nil, fiber.NewError(http.StatusUnauthorized, "Unauthorized"))
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestForgotPassword_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/forgot-password", handler.ForgotPassword)

	body := map[string]string{"email": "test@test.com"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("ForgotPassword", body["email"]).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/forgot-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSetPassword_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/set-password", handler.SetPassword)

	body := map[string]string{"token": "valid_token", "password": "new_password"}
	bodyJSON, _ := json.Marshal(body)
	mockService.On("SetPassword", body["token"], body["password"]).Return(nil)
	req := httptest.NewRequest(http.MethodPost, "/set-password", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateProfile_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Post("/profile", func(ctx *fiber.Ctx) error {
		ctx.Locals("userID", 1)
		return handler.CreateProfile(ctx)
	})

	body := dto.UserProfile{FirstName: "John", LastName: "Doe", Phone: "1234567890"}
	bodyJSON, _ := json.Marshal(body)

	mockService.On("CreateProfile", mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/profile", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid_token")
	resp, _ := app.Test(req)
	log.Printf("Response: %v", resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestGetProfile_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Get("/profile", func(ctx *fiber.Ctx) error {
		ctx.Locals("userID", 1)
		return handler.GetProfile(ctx)
	})

	mockService.On("GetProfile", 1).Return(&domain.User{ID: 1, Email: "test@test.com"}, nil)
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAuth_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Get("/auth", handler.Auth)

	mockService.On("Authenticate", mock.Anything).Return(&domain.User{ID: 1, Email: "test@test.com"}, nil)
	req := httptest.NewRequest(http.MethodGet, "/auth", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMe_Success(t *testing.T) {
	app, mockService, handler := setupTest(t)
	app.Get("/me", func(ctx *fiber.Ctx) error {
		ctx.Locals("userID", 1)
		return handler.Me(ctx)
	})
	mockService.On("GetProfile", 1).Return(&domain.User{
		ID:    1,
		Email: "test@test.com",
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
