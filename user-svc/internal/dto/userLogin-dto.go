package dto

type UserSignup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
