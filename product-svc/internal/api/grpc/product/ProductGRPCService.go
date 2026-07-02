package product

import (
	"context"
	"strconv"

	"github.com/sundayyogurt/product_service/internal/services"
)

type productGRPCService struct {
	UnimplementedProductServiceServer
	Service services.CatalogService
}

func (s *productGRPCService) GetProductById(ctx context.Context, req *GetProductRequest) (*ProductResponse, error) {

	product, err := s.Service.GetProductByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &ProductResponse{
		Id:         int32(product.ID),
		Name:       product.Name,
		CategoryId: strconv.Itoa(int(product.CategoryID)),
		Price:      product.Price,
		Stock:      int32(product.Stock),
	}, nil
}

func NewProductGRPCService(service services.CatalogService) *productGRPCService {
	return &productGRPCService{Service: service}
}
