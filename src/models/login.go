package models

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId uint64 `json:"userId"`
}
