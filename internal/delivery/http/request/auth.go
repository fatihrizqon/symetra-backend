package request

type RegisterRequest struct {
	Username string `validate:"required,min=1,max=20" json:"username" example:"johndoe"`
	Name     string `validate:"required,min=1,max=100" json:"name" example:"John Doe"`
	Email    string `validate:"required,email,min=1,max=254" json:"email" example:"[EMAIL_ADDRESS]"`
	Password string `validate:"required,min=8,max=100" json:"password" example:"yoursecretpassword"`
}

type LoginRequest struct {
	Email    string `validate:"required,email,min=1,max=254" json:"email" example:"[EMAIL_ADDRESS]"`
	Password string `validate:"required,min=8,max=100" json:"password" example:"yoursecretpassword"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}
