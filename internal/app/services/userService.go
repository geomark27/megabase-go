package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/database"
	"megabaseGo/internal/models"
	"megabaseGo/internal/utils"

	"gorm.io/gorm"
)

type UserService struct {
	hasher utils.PasswordHasher
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService() *UserService {
	return &UserService{
		hasher: utils.NewBcryptHasher(),
	}
}

func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	db := database.GetDB()

	// Verificar que el rol existe
	var role models.Role
	if err := db.First(&role, req.RoleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	// Verificar username único
	var existingUser models.User
	if err := db.Where("user_name = ?", req.UserName).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Verificar email único
	if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash de la contraseña
	hashedPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Valor por defecto para IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// 2. Generamos un token aleatorio para el campo RememberToken
	rememberToken, err := generateSecureToken(32)
	if err != nil {
		return nil, errors.New("failed to generate remember token")
	}

	// Crear usuario
	user := models.User{
		Name:     req.Name,
		UserName: req.UserName,
		Email:    req.Email,
		Password: hashedPassword,
		RoleID:   req.RoleID,
		IsActive: isActive,
		RememberToken: rememberToken,
	}

	// Guardar en BD
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Cargar relación y devolver
	if err := db.Preload("Role").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return s.toUserResponse(&user), nil
}

// GetUsers obtiene todos los usuarios con filtros opcionales
func (s *UserService) GetUsers(includeInactive bool, roleID *uint) ([]dto.UserResponse, error) {
	db := database.GetDB()
	var users []models.User

	query := db.Preload("Role")

	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}

	if roleID != nil {
		query = query.Where("role_id = ?", *roleID)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, *s.toUserResponse(&user))
	}

	return responses, nil
}

// GetUserByID obtiene un usuario por ID
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	db := database.GetDB()
	var user models.User

	if err := db.Preload("Role").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.toUserResponse(&user), nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	db := database.GetDB()
	var user models.User

	// Obtener usuario existente
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Verificar rol si se está cambiando
	if req.RoleID != 0 && req.RoleID != user.RoleID {
		var role models.Role
		if err := db.First(&role, req.RoleID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("role not found")
			}
			return nil, err
		}
	}

	// Verificar username único si se está cambiando
	if req.UserName != "" && req.UserName != user.UserName {
		var existing models.User
		if err := db.Where("user_name = ? AND id != ?", req.UserName, id).First(&existing).Error; err == nil {
			return nil, errors.New("username already exists")
		}
	}

	// Verificar email único si se está cambiando
	if req.Email != "" && req.Email != user.Email {
		var existing models.User
		if err := db.Where("email = ? AND id != ?", req.Email, id).First(&existing).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	// Actualizar campos
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.UserName != "" {
		user.UserName = req.UserName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.RoleID != 0 {
		user.RoleID = req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Hash nueva contraseña si se proporciona
	if req.Password != "" {
		hashedPassword, err := s.hasher.HashPassword(req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = hashedPassword
	}

	// Guardar cambios
	if err := db.Save(&user).Error; err != nil {
		return nil, err
	}

	// Cargar relación actualizada
	if err := db.Preload("Role").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return s.toUserResponse(&user), nil
}

// DeleteUser elimina un usuario (soft delete)
func (s *UserService) DeleteUser(id uint) error {
	db := database.GetDB()

	// Verificar que el usuario existe
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Soft delete
	return db.Delete(&user).Error
}

// CheckUsernameAvailability verifica si un username está disponible
func (s *UserService) CheckUsernameAvailability(username string) (bool, error) {
	var count int64
	
	// Contar cuántos usuarios tienen ese username
	err := database.DB.Model(&models.User{}).
		Where("user_name = ?", username).
		Count(&count).Error
		
	if err != nil {
		return false, utils.NewInternalServerError("Error verificando disponibilidad del username")
	}
	
	// Si count es 0, está disponible (true)
	// Si count > 0, no está disponible (false)
	return count == 0, nil
}

// CheckEmailAvailability verifica si un email está disponible
func (s *UserService) CheckEmailAvailability(email string) (bool, error) {
	var count int64
	
	// Contar cuántos usuarios tienen ese email
	err := database.DB.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
		
	if err != nil {
		return false, utils.NewInternalServerError("Error verificando disponibilidad del email")
	}
	
	// Si count es 0, está disponible (true)
	// Si count > 0, no está disponible (false)
	return count == 0, nil
}

// toUserResponse convierte un modelo User a UserResponse
func (s *UserService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		UserName:    user.UserName,
		Email:       user.Email,
		RoleID:      user.RoleID,
		Role:        user.Role,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}