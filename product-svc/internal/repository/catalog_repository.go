package repository

import (
	"github.com/sundayyogurt/product_service/internal/domain"
	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(category *domain.Category) error
	GetAllCategories() ([]domain.Category, error)
	GetCategoryById(id uint) (*domain.Category, error)
	UpdateCategory(id uint, update *domain.Category) error

	CreateProduct(product *domain.Product) error
	GetAllProducts() ([]domain.Product, error)
	GetProductById(id uint) (*domain.Product, error)
	UpdateProduct(id uint, update *domain.Product) error
}

func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &categoryRepository{db: db}
}

type categoryRepository struct {
	db *gorm.DB
}

func (c categoryRepository) CreateCategory(category *domain.Category) error {
	return c.db.Create(category).Error
}

func (c categoryRepository) GetAllCategories() ([]domain.Category, error) {
	var categories []domain.Category
	err := c.db.Find(&categories).Error
	return categories, err
}

func (c categoryRepository) GetCategoryById(id uint) (*domain.Category, error) {
	var category domain.Category
	err := c.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (c categoryRepository) UpdateCategory(id uint, update *domain.Category) error {
	return c.db.Model(&domain.Category{}).Where("id = ?", id).Updates(update).Error
}

func (c categoryRepository) CreateProduct(product *domain.Product) error {
	return c.db.Create(product).Error
}

func (c categoryRepository) GetAllProducts() ([]domain.Product, error) {
	var product []domain.Product
	err := c.db.Find(&product).Error
	return product, err
}

func (c categoryRepository) GetProductById(id uint) (*domain.Product, error) {
	var product domain.Product
	err := c.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (c categoryRepository) UpdateProduct(id uint, update *domain.Product) error {
	return c.db.Model(&domain.Product{}).Where("id = ?", id).Updates(update).Error
}
