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

// NewUserService dependency injection
func NewUserService(repo repository.UserRepository, producer interfaces.ProducerHandler) UserService {
	return &userService{repo: repo, producer: producer}
}

type userService struct {
	repo     repository.UserRepository
	producer interfaces.ProducerHandler
}

func (s *userService) Register(input dto.UserSignup) error {
	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error on hashing password: %v", err)
		return errors.New("failed to hash password")
	}

	// save to DB
	err = s.repo.CreateUser(&domain.User{
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

func (s *userService) Login(email, password string) (*domain.User, error) {
	// find the existing user
	user, err := s.repo.FindUserByEmail(email)

	if user == nil || err != nil {
		log.Printf("FindUserByEmail: %v", email)
		return nil, errors.New("invalid email or password")
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *userService) ForgotPassword(email string) error {
	// find existing user
	user, err := s.repo.FindUserByEmail(email)

	if user == nil || err != nil {
		log.Printf("FindUserByEmail: %v", email)
		return errors.New("invalid email or password")
	}

	// generate unique reset token
	resetToken, err := bcrypt.GenerateFromPassword([]byte(user.Email), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to generate reset token: %v", err)
		return errors.New("failed to generate reset token")
	}

	user.ResetToken = string(resetToken)
	err = s.repo.SaveUser(user)

	if err != nil {
		log.Printf("failed to save user: %v", err)
		return errors.New("failed to save user")
	}

	// publish reset token to Kafka or send email
	// s.producer.PublishMessage([]byte("user.reset_password"), []byte(resetToken))

	return nil
}

func (s *userService) SetPassword(token, password string) error {
	user, err := s.repo.FindUserByResetToken(token)
	if err != nil {
		log.Printf("invalid or expired token: %v", token)
		return errors.New("invalid or expired token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("failed to hash password: %v", err)
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	user.ResetToken = ""

	err = s.repo.SaveUser(user)
	if err != nil {
		log.Printf("failed to save user: %v", err)
		return errors.New("failed to save user")
	}

	return nil
}

func (s *userService) CreateProfile(profile dto.UserProfile) error {

	user, err := s.repo.FindUserById(profile.UserID)
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return errors.New("failed to find user")
	}

	if user.FirstName != "" {
		return errors.New("first name is already taken")
	}

	address := domain.Address{
		AddressLine1: profile.Address.AddressLine1,
		AddressLine2: profile.Address.AddressLine2,
		City:         profile.Address.City,
		Country:      profile.Address.Country,
		PostCode:     profile.Address.PostCode,
	}

	user.Address = address
	err = s.repo.SaveUser(user)
	if err != nil {
		log.Printf("failed to save user: %v", err)
		return errors.New("failed to save user")
	}

	return nil
}

func (s *userService) GetProfile(userID int) (*domain.User, error) {
	profile, err := s.repo.FindUserById(userID)
	if err != nil {
		log.Printf("failed to find user: %v", err)
		return nil, errors.New("failed to find user")
	}

	return profile, nil
}

func (s *userService) Authenticate(c *fiber.Ctx) (*domain.User, error) {
	user := c.Locals("userID")

	authUser, err := s.repo.FindUserById(user.(int))
	if err != nil {
		log.Printf("failed tio find user: %v", err)
		return nil, errors.New("failed to find user")
	}
	return authUser, nil
}
