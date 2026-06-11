package dto

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	ImageURL    string  `json:"image_url"`
	CategoryID  uint    `json:"category_id"`
	Status      string  `json:"status"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       uint    `json:"stock"`
	ImageURL    string  `json:"image_url"`
	Status      string  `json:"status"`
}
