package repository

import (
	"errors"

	"github.com/sundayyogurt/order_service/internal/domain"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetOrCreate(userID uint) (*domain.Cart, error)
	Get(userID uint) (*domain.Cart, error)
	AddOrIncrement(userID uint, item domain.CartItem) error
	UpdateQty(userID uint, productID uint, qty uint) error
	Remove(userID uint, productID uint) error
	Clear(userID uint) error
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

type cartRepository struct{ db *gorm.DB }

func (c cartRepository) GetOrCreate(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	err := c.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = domain.Cart{UserID: userID}
			if err := c.db.Create(&cart).Error; err != nil {
				return nil, err
			}
			return &cart, nil
		}
		return nil, err
	}
	return &cart, nil
}

func (c cartRepository) Get(userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	if err := c.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (c cartRepository) AddOrIncrement(userID uint, item domain.CartItem) error {
	cart, err := c.Get(userID)
	if err != nil {
		return err
	}

	var existing domain.CartItem
	err = c.db.Where("cart_id = ? AND product_id = ?", cart.ID, item.ProductID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item.CardID = cart.ID
		return c.db.Create(&item).Error
	}

	if err != nil {
		return err
	}

	newQty := existing.Qty + item.Qty
	return c.db.Model(&existing).Updates(map[string]interface{}{
		"qty":          newQty,
		"price":        item.Price,
		"product_name": item.ProductName,
		"image_url":    item.ImageURL,
	}).Error
}

func (c cartRepository) UpdateQty(userID uint, productID uint, qty uint) error {
	cart, err := c.GetOrCreate(userID)
	if err != nil {
		return err
	}

	if qty == 0 {
		return c.Remove(userID, productID)
	}

	return c.db.Model(&domain.CartItem{}).Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		Updates(map[string]interface{}{
			"qty": qty,
		}).Error

}

func (c cartRepository) Remove(userID uint, productID uint) error {
	cart, err := c.GetOrCreate(userID)
	if err != nil {
		return err
	}

	return c.db.Where("cart_id = ? AND product_id = ?", cart.ID, productID).
		Delete(&domain.CartItem{}).Error

}

func (c cartRepository) Clear(userID uint) error {
	cart, err := c.GetOrCreate(userID)
	if err != nil {
		return err
	}

	return c.db.Where("cart_id = ?", cart.ID).Delete(&domain.CartItem{}).Error
}
