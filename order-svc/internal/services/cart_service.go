package services

import (
	"errors"

	"github.com/sundayyogurt/order_service/internal/domain"
	"github.com/sundayyogurt/order_service/internal/dto"
	"github.com/sundayyogurt/order_service/internal/repository"
)

type CartService interface {
	Add(userID uint, in dto.CartAddRequest) error
	UpdateQty(userID uint, productID, qty uint) error
	Remove(userID uint, productID uint) error
	Get(userID uint) (*dto.CartResponse, error)
}

type cartService struct {
	repo repository.CartRepository
}

func NewCartService(r repository.CartRepository) CartService {
	return &cartService{repo: r}
}

func (s *cartService) Add(userID uint, in dto.CartAddRequest) error {
	if in.Qty == 0 {
		return errors.New("qty must be greater than zero")
	}

	return s.repo.AddOrIncrement(userID, domain.CartItem{
		ProductID:   in.ProductID,
		ProductName: in.ProductName,
		ImageURL:    in.ImageURL,
		Qty:         in.Qty,
		Price:       in.Price,
	})
}

func (s *cartService) UpdateQty(userID, productID uint, qty uint) error {
	return s.repo.UpdateQty(userID, productID, qty)
}

func (s *cartService) Remove(userID uint, productID uint) error {
	return s.repo.Remove(userID, productID)
}

func (s *cartService) Get(userID uint) (*dto.CartResponse, error) {
	c, err := s.repo.GetOrCreate(userID)
	if err != nil {
		return nil, err
	}

	resp := dto.CartResponse{UserID: userID}
	var subtotal float64

	for _, item := range c.Items {
		line := float64(item.Qty) * item.Price
		subtotal += line

		resp.Items = append(resp.Items, dto.CardItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ImageURL:    item.ImageURL,
			Qty:         item.Qty,
			Price:       item.Price,
			LineTotal:   line,
		})
	}
	resp.Count = len(resp.Items)
	resp.Subtotal = subtotal

	return &resp, nil
}
