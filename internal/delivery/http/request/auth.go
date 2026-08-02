package request

type RegisterRequest struct {
	Username string `validate:"required,min=3,max=20,alphanum" json:"username"`
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required,email,min=1,max=254" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}
