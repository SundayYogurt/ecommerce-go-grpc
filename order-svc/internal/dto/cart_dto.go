package dto

type CartAddRequest struct {
	ProductID   uint    `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name"`
	ImageURL    string  `json:"image_url"`
	Qty         uint    `json:"qty" validate:"required, gt=0"`
	Price       float64 `json:"price" validate:"required,gt=0"`
}

type CardUpdateQtyRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Qty       uint `json:"qty" validate:"required,gte=0"`
}
type CartRemoveRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
}

type CardItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	ImageURL    string  `json:"image_url"`
	Qty         uint    `json:"qty"`
	Price       float64 `json:"price"`
	LineTotal   float64 `json:"line_total"`
}

type CartResponse struct {
	UserID   uint               `json:"user_id"`
	Items    []CardItemResponse `json:"items"`
	Count    int                `json:"count"`
	Subtotal float64            `json:"subtotal"`
}
