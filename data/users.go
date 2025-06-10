package data

type CreateUsersRequest struct {
	Username   string `validate:"required,min=1,max=255" json:"username"`
	Name       string `validate:"required,min=1,max=255" json:"name"`
	Password   string `validate:"required,min=8,max=255" json:"password"`
	Email      string `validate:"required,email" json:"email"`
	AvatarType string `validate:"omitempty" json:"avatar_type"`
}

type UserResponse struct {
	Id       string    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
