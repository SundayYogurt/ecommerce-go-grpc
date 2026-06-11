package domain

import "gorm.io/gorm"

type Product struct {
	Id          int     `json:"id" gorm:"primary_key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	ImageURL    string  `json:"image_url"`
	Status      string  `json:"status"` // publis or unpublish
	gorm.Model
}
