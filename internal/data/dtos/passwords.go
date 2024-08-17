package dtos

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordReset struct {
	Token    string `json:"token"    validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
