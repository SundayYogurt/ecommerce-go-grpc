package domain

import "gorm.io/gorm"

type Address struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	City         string `json:"city"`
	PostCode     string `json:"post_code"`
	Country      string `json:"country"`
	UserId       uint   `json:"user_id"`
	gorm.Model
}
