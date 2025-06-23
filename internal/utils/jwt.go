package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims estructura para los claims del JWT
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

// JWTManager maneja la generación y validación de tokens JWT
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// NewJWTManager crea una nueva instancia del manager JWT
func NewJWTManager() *JWTManager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "megabase-default-secret-key-change-in-production"
	}

	duration := time.Hour * 24 // 24 horas por defecto
	if durationStr := os.Getenv("JWT_DURATION_HOURS"); durationStr != "" {
		if hours, err := strconv.Atoi(durationStr); err == nil {
			duration = time.Hour * time.Duration(hours)
		}
	}

	return &JWTManager{
		secretKey:     secret,
		tokenDuration: duration,
	}
}

// GenerateToken genera un nuevo token JWT
func (manager *JWTManager) GenerateToken(userID uint, userName, email string, roleID uint, roleName string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		RoleID:   roleID,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "megabase-go",
			Subject:   strconv.Itoa(int(userID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// GenerateRefreshToken genera un refresh token (válido por más tiempo)
func (manager *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 7 días
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "megabase-go-refresh",
		Subject:   strconv.Itoa(int(userID)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// ValidateToken valida un token JWT y retorna los claims
func (manager *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateRefreshToken valida un refresh token
func (manager *JWTManager) ValidateRefreshToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(manager.secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.Issuer != "megabase-go-refresh" {
			return 0, errors.New("invalid refresh token issuer")
		}
		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, err
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid refresh token")
}

// GetTokenDuration retorna la duración del token en segundos
func (manager *JWTManager) GetTokenDuration() int64 {
	return int64(manager.tokenDuration.Seconds())
}