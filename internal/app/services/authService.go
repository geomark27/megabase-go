package services

import (
	"errors"
	"time"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/database"
	"megabaseGo/internal/models"
	"megabaseGo/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	userService *UserService
	jwtManager  *utils.JWTManager
	hasher      utils.PasswordHasher
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService() *AuthService {
	return &AuthService{
		userService: NewUserService(),
		jwtManager:  utils.NewJWTManager(),
		hasher:      utils.NewBcryptHasher(),
	}
}

// Login autentica un usuario y retorna tokens
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	db := database.GetDB()

	// Buscar usuario por username con rol
	var user models.User
	if err := db.Preload("Role").Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Verificar que el usuario esté activo
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Verificar contraseña
	if err := s.hasher.ComparePassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Actualizar último login
	user.LastLoginAt = time.Now()
	db.Save(&user)

	// Generar tokens
	accessToken, err := s.jwtManager.GenerateToken(
		user.ID,
		user.UserName,
		user.Email,
		user.RoleID,
		user.Role.Name,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Crear respuesta
	userResponse := s.toUserResponse(&user)

	return &dto.AuthResponse{
		User:         *userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetTokenDuration(),
	}, nil
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Usar el UserService para crear el usuario
	createUserReq := &dto.CreateUserRequest{
		Name:     req.Name,
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   req.RoleID,
	}

	userResponse, err := s.userService.CreateUser(createUserReq)
	if err != nil {
		return nil, err
	}

	// Obtener el usuario completo con rol para generar tokens
	db := database.GetDB()
	var user models.User
	if err := db.Preload("Role").First(&user, userResponse.ID).Error; err != nil {
		return nil, err
	}

	// Generar tokens
	accessToken, err := s.jwtManager.GenerateToken(
		user.ID,
		user.UserName,
		user.Email,
		user.RoleID,
		user.Role.Name,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &dto.AuthResponse{
		User:         *userResponse,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetTokenDuration(),
	}, nil
}

// RefreshToken genera un nuevo access token usando el refresh token
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	// Validar refresh token
	userID, err := s.jwtManager.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Obtener usuario actualizado
	db := database.GetDB()
	var user models.User
	if err := db.Preload("Role").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Verificar que el usuario siga activo
	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	// Generar nuevos tokens
	accessToken, err := s.jwtManager.GenerateToken(
		user.ID,
		user.UserName,
		user.Email,
		user.RoleID,
		user.Role.Name,
	)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	userResponse := s.toUserResponse(&user)

	return &dto.AuthResponse{
		User:         *userResponse,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetTokenDuration(),
	}, nil
}

// ChangePassword cambia la contraseña de un usuario autenticado
func (s *AuthService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	db := database.GetDB()

	// Obtener usuario
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Verificar contraseña actual
	if err := s.hasher.ComparePassword(user.Password, req.CurrentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Hash nueva contraseña
	hashedPassword, err := s.hasher.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// Actualizar contraseña
	user.Password = hashedPassword
	if err := db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

// GetCurrentUser obtiene la información del usuario actual
func (s *AuthService) GetCurrentUser(userID uint) (*dto.UserResponse, error) {
	return s.userService.GetUserByID(userID)
}

// ValidateToken valida un token y retorna los claims
func (s *AuthService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtManager.ValidateToken(tokenString)
}

// toUserResponse convierte un modelo User a UserResponse
func (s *AuthService) toUserResponse(user *models.User) *dto.UserResponse {
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