package data

type CreateUsersRequest struct {
	Username string `validate:"required,min=1,max=255" json:"username"`
	Email    string `validate:"required,email" json:"email"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
