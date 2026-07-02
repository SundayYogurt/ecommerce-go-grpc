package domain

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	CardID      uint    `json:"card_id" gorm:"index"`
	ProductID   uint    `json:"product_id" gorm:"index"`
	ProductName string  `json:"product_name"`
	ImageURL    string  `json:"image_url"`
	Qty         uint    `json:"qty"`
	Price       float64 `json:"price"`
}

type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id" gorm:"uniqueIndex"`
	Items  []CartItem `json:"items" gorm:"constraint:OnDelete:CASCADE"`
}

func (CartItem) TableName() string { return "cart_items" }
