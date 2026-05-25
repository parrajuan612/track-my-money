package domain

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}
