package dto

import "github.com/google/uuid"

type CreateUsersRequest struct {
	Username   string `validate:"required,min=1,max=255" json:"username"`
	Name       string `validate:"required,min=1,max=255" json:"name"`
	Password   string `validate:"required,min=8,max=255" json:"password"`
	Email      string `validate:"required,email" json:"email"`
	AvatarType string `validate:"omitempty" json:"avatar_type"`
}

type UserResponse struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	Username   string    `json:"username"`
	AvatarType string    `json:"avatar_type,omitempty"`
	Password   string    `json:"password,omitempty"`
}

type UpdateAvatarUserRequest struct {
	Id         uuid.UUID `validate:"required" json:"id"`
	AvatarType string `validate:"required" json:"avatar_type"`
}
