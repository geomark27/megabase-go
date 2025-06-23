package dto

import "megabaseGo/internal/models"

// CreateUserRequest estructura para crear un usuario
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   uint   `json:"role_id" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

// UpdateUserRequest estructura para actualizar un usuario
type UpdateUserRequest struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	RoleID   uint   `json:"role_id"`
	IsActive *bool  `json:"is_active"`
}

// UserResponse estructura para respuestas (sin contraseña)
type UserResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	UserName    string      `json:"user_name"`
	Email       string      `json:"email"`
	RoleID      uint        `json:"role_id"`
	Role        models.Role `json:"role"`
	IsActive    bool        `json:"is_active"`
	LastLoginAt interface{} `json:"last_login_at"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
}