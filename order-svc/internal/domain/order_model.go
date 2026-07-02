package domain

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID     uint         `json:"user_id"`
	Total      float64      `json:"total"`
	Status     string       `json:"status"`
	OrderRef   string       `json:"order_ref"`
	OrderItems []OrderItems `json:"order_items" gorm:"constraint:OnDelete:CASCADE"`
}

type OrderItems struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  uint    `json:"quantity"`
	Price     float64 `json:"price"`
}
