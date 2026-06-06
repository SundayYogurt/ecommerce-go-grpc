package domain

import "gorm.io/gorm"

type User struct {
	ID         uint    `gorm:"primary_key" json:"id"`
	Email      string  `gorm:"uniqueIndex" json:"email"`
	Password   string  `json:"-"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Phone      string  `json:"phone"`
	ResetToken string  `json:"reset_token"`
	Address    Address `json:"address"` // relationship to user model
	gorm.Model
}
