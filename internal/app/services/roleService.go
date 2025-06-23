package services

import (
	"errors"
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/database"
	"megabaseGo/internal/models"

	"gorm.io/gorm"
)

type RoleService struct{}

// NewRoleService crea una nueva instancia del servicio de roles
func NewRoleService() *RoleService {
	return &RoleService{}
}

// CreateRole crea un nuevo rol
func (s *RoleService) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	db := database.GetDB()

	// Verificar nombre único
	var existingRole models.Role
	if err := db.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, errors.New("role with this name already exists")
	}

	// Valor por defecto para IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Crear rol
	role := models.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		IsActive:    isActive,
	}

	// Guardar en BD
	if err := db.Create(&role).Error; err != nil {
		return nil, err
	}

	return s.toRoleResponse(&role), nil
}

// GetRoles obtiene todos los roles con filtros opcionales
func (s *RoleService) GetRoles(includeInactive bool) ([]dto.RoleResponse, error) {
	db := database.GetDB()
	var roles []models.Role

	query := db
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	var responses []dto.RoleResponse
	for _, role := range roles {
		responses = append(responses, *s.toRoleResponse(&role))
	}

	return responses, nil
}

// GetRoleByID obtiene un rol por ID
func (s *RoleService) GetRoleByID(id uint) (*dto.RoleResponse, error) {
	db := database.GetDB()
	var role models.Role

	if err := db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return s.toRoleResponse(&role), nil
}

// UpdateRole actualiza un rol existente
func (s *RoleService) UpdateRole(id uint, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	db := database.GetDB()
	var role models.Role

	// Obtener rol existente
	if err := db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// Verificar nombre único si se está cambiando
	if req.Name != "" && req.Name != role.Name {
		var existing models.Role
		if err := db.Where("name = ? AND id != ?", req.Name, id).First(&existing).Error; err == nil {
			return nil, errors.New("role with this name already exists")
		}
	}

	// Actualizar campos
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	// Guardar cambios
	if err := db.Save(&role).Error; err != nil {
		return nil, err
	}

	return s.toRoleResponse(&role), nil
}

// DeleteRole elimina un rol (soft delete)
func (s *RoleService) DeleteRole(id uint) error {
	db := database.GetDB()

	// Verificar que el rol existe
	var role models.Role
	if err := db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return err
	}

	// Verificar que no hay usuarios usando este rol
	var userCount int64
	if err := db.Model(&models.User{}).Where("role_id = ?", id).Count(&userCount).Error; err != nil {
		return err
	}

	if userCount > 0 {
		return errors.New("cannot delete role: it is assigned to users")
	}

	// Soft delete
	return db.Delete(&role).Error
}

// toRoleResponse convierte un modelo Role a RoleResponse
func (s *RoleService) toRoleResponse(role *models.Role) *dto.RoleResponse {
	return &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: role.Description,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}