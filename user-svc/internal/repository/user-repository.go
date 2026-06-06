package repository

import (
	"github.com/sundayyogurt/user_service/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	FindUserByEmail(email string) (*domain.User, error)
	SaveUser(user *domain.User) error
	FindUserByResetToken(token string) (*domain.User, error)
	FindUserById(userID int) (*domain.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

type userRepository struct {
	db *gorm.DB
}

func (r *userRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Preload("Address").First(&user, "email=?", email).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveUser(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) FindUserByResetToken(token string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("reset_token=?", token).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserById(userID int) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
