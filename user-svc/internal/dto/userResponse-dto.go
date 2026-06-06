package dto

type AuthResponse struct {
	UserID int     `json:"user_id"`
	Email  string  `json:"email"`
	Exp    float64 `json:"exp"`
}
