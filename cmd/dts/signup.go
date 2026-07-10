package dts

type UserSignUp struct {
	UserName  string `json:"username" validate:"required,min=5,max=50"`
	FirstName string `json:"first_name" validate:"required,min=3,max=25"`
	LastName  string `json:"last_name" validate:"required,min=3,max=25"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=5,max=72"`
}
