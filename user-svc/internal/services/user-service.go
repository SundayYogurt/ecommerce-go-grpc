package services

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sundayyogurt/user_service/internal/domain"
	"github.com/sundayyogurt/user_service/internal/dto"
	"github.com/sundayyogurt/user_service/internal/interfaces"
	"github.com/sundayyogurt/user_service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(input dto.UserSignup) error
	Login(email, password string) (*domain.User, error)
	ForgotPassword(email string) error
	SetPassword(token, password string) error
	CreateProfile(profile dto.UserProfile) error
	GetProfile(userID int) (*domain.User, error)
	Authenticate(c *fiber.Ctx) (*domain.User, error)
}

func NewUserService(repo repository.UserRepository, producer *interfaces.ProducerHandler) UserService {
	return &userService{repo: repo, producer: producer}
}

type userService struct {
	repo     repository.UserRepository
	producer *interfaces.ProducerHandler
}

func (u userService) Register(input dto.UserSignup) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error on hashing password: %v", err)
		return errors.New("failed to hash password")
	}

	// save to DB
	err = u.repo.CreateUser(&domain.User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Phone:    input.Phone,
	})

	// publish to Kafka to send the confirmation email
	// s.producer.PublishMessage([]byte("user.registered", []byte(input.Email)))
	if err != nil {
		log.Printf("Error on saving user: %v", err)
		return errors.New("failed to create user")
	}

	// return err or nil
	return nil
}

func (u userService) Login(email, password string) (*domain.User, error) {
	return nil, nil
}

func (u userService) ForgotPassword(email string) error {
	return nil
}

func (u userService) SetPassword(token, password string) error {
	return nil
}

func (u userService) CreateProfile(profile dto.UserProfile) error {
	return nil
}

func (u userService) GetProfile(userID int) (*domain.User, error) {
	return nil, nil
}

func (u userService) Authenticate(c *fiber.Ctx) (*domain.User, error) {
	return nil, nil
}
