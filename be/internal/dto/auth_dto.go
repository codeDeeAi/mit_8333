package dto

import "time"

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	FullName string  `json:"full_name" binding:"required,min=2,max=150"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6,max=72"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,max=30"`
	RoleID   int64   `json:"role_id" binding:"omitempty,gt=0"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// UserResponse is the public user payload returned by auth endpoints.
type UserResponse struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	Role      string    `json:"role,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthResponse is returned for login and register endpoints.
type AuthResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// RoleResponse is a single role exposed to clients.
type RoleResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// RegistrationDataResponse holds reference data for the sign-up screen.
type RegistrationDataResponse struct {
	Roles []RoleResponse `json:"roles"`
}
