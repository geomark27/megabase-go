package handlers

import (
	"net/http"
	"time"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/middleware"
	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	authResponse, err := h.authService.Login(&req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	accessTokenMaxAge := int(authResponse.ExpiresIn)
	c.SetCookie("access_token", authResponse.AccessToken, accessTokenMaxAge, "/", "localhost", false, true)
	refreshTokenMaxAge := int((time.Hour * 24 * 7).Seconds())
	c.SetCookie("refresh_token", authResponse.RefreshToken, refreshTokenMaxAge, "/", "localhost", false, true)

	utils.SendSuccess(c, http.StatusOK, "Login exitoso", gin.H{"user": authResponse.User})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	accessTokenMaxAge := int(authResponse.ExpiresIn)
	c.SetCookie("access_token", authResponse.AccessToken, accessTokenMaxAge, "/", "localhost", false, true)
	refreshTokenMaxAge := int((time.Hour * 24 * 7).Seconds())
	c.SetCookie("refresh_token", authResponse.RefreshToken, refreshTokenMaxAge, "/", "localhost", false, true)

	utils.SendSuccess(c, http.StatusCreated, "Registro exitoso", gin.H{"user": authResponse.User})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		utils.HandleGinError(c, utils.NewUnauthorizedError("Refresh token no encontrado"))
		return
	}

	authResponse, err := h.authService.RefreshToken(&dto.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	accessTokenMaxAge := int(authResponse.ExpiresIn)
	c.SetCookie("access_token", authResponse.AccessToken, accessTokenMaxAge, "/", "localhost", false, true)
	refreshTokenMaxAge := int((time.Hour * 24 * 7).Seconds())
	c.SetCookie("refresh_token", authResponse.RefreshToken, refreshTokenMaxAge, "/", "localhost", false, true)

	utils.SendSuccess(c, http.StatusOK, "Token refrescado exitosamente", gin.H{"user": authResponse.User})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	utils.SendSuccess(c, http.StatusOK, "Logout exitoso", nil)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Perfil obtenido exitosamente", gin.H{"user": user})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := middleware.GetCurrentUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "Usuario no autenticado")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	err := h.authService.ChangePassword(userID, &req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Contraseña cambiada exitosamente", nil)
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
	claims, exists := middleware.GetCurrentUserClaims(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "No autenticado")
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Usuario autenticado", gin.H{
		"authenticated": true,
		"user": gin.H{
			"id":        claims.UserID,
			"user_name": claims.UserName,
			"email":     claims.Email,
			"role_id":   claims.RoleID,
			"role_name": claims.RoleName,
		},
	})
}