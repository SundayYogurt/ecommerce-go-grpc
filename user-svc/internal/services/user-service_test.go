package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sundayyogurt/user_service/internal/domain"
	"github.com/sundayyogurt/user_service/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) FindUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}

	return user.(*domain.User), args.Error(1)
}

func (m *mockUserRepository) SaveUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepository) FindUserByResetToken(token string) (*domain.User, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepository) FindUserById(userID int) (*domain.User, error) {
	args := m.Called(userID)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}

	return user.(*domain.User), args.Error(1)
}

type MockProducerHandler struct {
	mock.Mock
}

func (m *MockProducerHandler) PublishMessage(key, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {

	// mock dependencies that need to pass to mimic actual behavior
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)

	// create a new user service with the mocked dependencies
	svc := NewUserService(repo, producer)

	// create a new user signup input
	input := dto.UserSignup{
		Email:    "test@test.com",
		Password: "password",
		Phone:    "0812345678",
	}

	// set up the expected for the CreateUser method to mimic actual behavior
	repo.On("CreateUser", mock.MatchedBy(func(user *domain.User) bool {
		assert.Equal(t, input.Email, user.Email)
		assert.Equal(t, input.Phone, user.Phone)
		assert.Len(t, user.Password, 60) // bcrypt hash lenght
		return true
	})).Return(nil)

	// call the Register method of user service
	err := svc.Register(input)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)
	svc := NewUserService(repo, producer)

	email := "test@test.com"
	password := "password"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	repo.On("FindUserByEmail", email).Return(user, nil)

	result, err := svc.Login(email, password)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.Password, result.Password)

	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)
	svc := NewUserService(repo, producer)

	email := "test@test.com"
	password := "wrong_password"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("right_password"), bcrypt.DefaultCost)

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	repo.On("FindUserByEmail", email).Return(user, nil)

	result, err := svc.Login(email, password)
	assert.Nil(t, result) // no result
	assert.Error(t, err, "invalid email or password")

	repo.AssertExpectations(t)
}

func TestForgotPassword_UserNotFound(t *testing.T) {
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)
	svc := NewUserService(repo, producer)

	email := "test@notest.com"
	repo.On("FindUserByEmail", email).Return(nil, errors.New("user not found"))
	err := svc.ForgotPassword(email)
	assert.Error(t, err, "user not found")
}

func TestSetPassword_InvalidToken(t *testing.T) {
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)
	svc := NewUserService(repo, producer)

	repo.On("FindUserByResetToken", "invalid_token").Return((*domain.User)(nil), errors.New("not found"))

	err := svc.SetPassword("invalid_token", "new_password")

	assert.Error(t, err, "invalid or expired token")
	repo.AssertExpectations(t)
}

func TestSetPassword_Success(t *testing.T) {
	repo := new(mockUserRepository)
	producer := new(MockProducerHandler)
	svc := NewUserService(repo, producer)

	repo.On("FindUserByResetToken", "valid_token").Return(&domain.User{
		ID:    1,
		Email: "test@test.com",
	}, nil)

	repo.On("SaveUser", mock.MatchedBy(func(u *domain.User) bool {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte("new_password"))
		return u.Email == "test@test.com" && err == nil
	})).Return(nil)

	err := svc.SetPassword("valid_token", "new_password")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestGetProfile_Success(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo, nil)

	userID := 1
	user := &domain.User{
		ID: uint(userID),
	}
	repo.On("FindUserById", userID).Return(user, nil)

	profile, err := svc.GetProfile(userID)
	assert.NoError(t, err)
	assert.Equal(t, user, profile)
	repo.AssertExpectations(t)
}

func TestGetProfile_NotFound(t *testing.T) {
	repo := new(mockUserRepository)
	svc := NewUserService(repo, nil)

	userID := 1
	repo.On("FindUserById", userID).Return((*domain.User)(nil), errors.New("not found"))

	profile, err := svc.GetProfile(userID)
	assert.Nil(t, profile)
	assert.Error(t, err, "profile not found")
	repo.AssertExpectations(t)
}
