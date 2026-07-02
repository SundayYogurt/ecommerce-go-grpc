package dto

type CreateOrderRequest struct {
	UserID uint `json:"user_id"`
	Items  []struct {
		ProductID uint `json:"product_id"`
		Quantity  uint `json:"quantity"`
		Price     uint `json:"price"`
	} `json:"items"`
}

type CreateOrderResponse struct {
	OrderID  uint   `json:"order_id"`
	OrderRef string `json:"order_ref"`
	Total    uint   `json:"total"`
	Status   string `json:"status"`
}
