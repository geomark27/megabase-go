package dto

// CreateRoleRequest estructura para crear un rol
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// UpdateRoleRequest estructura para actualizar un rol
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

// RoleResponse estructura para respuestas
type RoleResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name"`
	Description string      `json:"description"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   interface{} `json:"created_at"`
	UpdatedAt   interface{} `json:"updated_at"`
}