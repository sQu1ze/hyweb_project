package model

import "time"

type User struct {
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=50" example:"admin2@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"admin@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"password123"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"New@password123"`
}
