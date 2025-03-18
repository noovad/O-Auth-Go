package data

type CreateUsersRequest struct {
	Username string `validate:"required,min=1,max=200" json:"username"`
	Email    string `validate:"required" json:"email"`
}

type UsersResponse struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}
