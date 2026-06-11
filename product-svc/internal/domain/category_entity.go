package domain

import "gorm.io/gorm"

type Category struct {
	Id          int    `json:"id" gorm:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Status      string `json:"status"`
	gorm.Model
}
