package services

import (
	"errors"
	"log"

	"github.com/sundayyogurt/product_service/internal/domain"
	"github.com/sundayyogurt/product_service/internal/dto"
	"github.com/sundayyogurt/product_service/internal/interfaces"
	"github.com/sundayyogurt/product_service/internal/repository"
)

type CatalogService interface {
	CreateCategory(input dto.CreateCategoryRequest) error
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id uint) (*domain.Category, error)
	UpdateCategory(id uint, input dto.UpdateCategoryRequest) error

	CreateProduct(input dto.CreateProductRequest) error
	GetAllProducts() ([]domain.Product, error)
	GetProductByID(id uint) (*domain.Product, error)
	UpdateProduct(id uint, input dto.UpdateProductRequest) error
	interfaces.ConsumerHandler
}

func NewCatalogService(repo repository.CatalogRepository) CatalogService {
	return &catalogService{repo: repo}
}

type catalogService struct {
	repo repository.CatalogRepository
}

func (c catalogService) HandleMessage(message string) error {
	log.Println("Received kafka messages inside product service:", message)
	return nil
}

func (c catalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	category := domain.Category{
		Name:        input.Name,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Status:      "publish",
	}

	err := c.repo.CreateCategory(&category)
	if err != nil {
		// Log the error to push to grafana to visualize the logs
		log.Printf("Error creating category: %v", err)
		return errors.New("error on create category") // integrate Grafana logging here
	}

	return nil
}

func (c catalogService) GetAllCategories() ([]domain.Category, error) {
	return c.repo.GetAllCategories()
}

func (c catalogService) GetCategoryByID(id uint) (*domain.Category, error) {
	return c.repo.GetCategoryById(id)
}

func (c catalogService) UpdateCategory(id uint, input dto.UpdateCategoryRequest) error {
	existingCategory, err := c.repo.GetCategoryById(id)
	if err != nil {
		log.Printf("Error getting category by ID: %v", err)
		return errors.New("category not found")
	}

	if input.Name != "" {
		existingCategory.Name = input.Name
	}

	if input.Description != "" {
		existingCategory.Description = input.Description
	}

	if input.Status != "" {
		existingCategory.Status = input.Status
	}

	return c.repo.UpdateCategory(id, existingCategory)

}

func (c catalogService) CreateProduct(input dto.CreateProductRequest) error {
	product := domain.Product{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.CategoryID,
		Price:       input.Price,
		Stock:       input.Stock,
		ImageURL:    input.ImageURL,
		Status:      "publish",
	}

	err := c.repo.CreateProduct(&product)
	if err != nil {
		// Log the error to push to grafana to visualize the logs
		log.Printf("Error creating product: %v", err)
		return errors.New("error on create product") // integrate Grafana logging here
	}

	return nil
}

func (c catalogService) GetAllProducts() ([]domain.Product, error) {
	return c.repo.GetAllProducts()
}

func (c catalogService) GetProductByID(id uint) (*domain.Product, error) {
	return c.repo.GetProductById(id)
}

func (c catalogService) UpdateProduct(id uint, input dto.UpdateProductRequest) error {
	existingProduct, err := c.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error getting product by ID: %v", err)
		return errors.New("product not found")
	}

	if input.Name != "" {
		existingProduct.Name = input.Name
	}

	if input.Description != "" {
		existingProduct.Description = input.Description
	}

	if input.Price != 0 {
		existingProduct.Price = input.Price
	}

	if input.Stock != 0 {
		existingProduct.Stock = input.Stock
	}

	if input.ImageURL != "" {
		existingProduct.ImageURL = input.ImageURL
	}

	if input.Status != "" {
		existingProduct.Status = input.Status
	}

	return c.repo.UpdateProduct(id, existingProduct)
}
