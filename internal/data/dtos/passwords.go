package dtos

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordReset struct {
	Token    string `json:"token"         validate:"required"`
	Password string `json:"new_password1" validate:"required,min=8"`
}
