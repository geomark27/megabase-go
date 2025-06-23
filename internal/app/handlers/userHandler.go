package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(),
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	// CREATE: Mantiene mensaje
	utils.SendSuccess(c, http.StatusCreated, "Usuario creado correctamente", gin.H{"user": user})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"
	var roleID *uint
	if roleIDStr := c.Query("role_id"); roleIDStr != "" {
		if id, err := strconv.ParseUint(roleIDStr, 10, 32); err == nil {
			roleIDUint := uint(id)
			roleID = &roleIDUint
		}
	}

	users, err := h.userService.GetUsers(includeInactive, roleID)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	// GET: Solo status y data (sin mensaje)
	utils.SendData(c, http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de usuario inválido"))
		return
	}

	user, err := h.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	// GET: Solo status y data (sin mensaje)
	utils.SendData(c, http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de usuario inválido"))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	user, err := h.userService.UpdateUser(uint(userID), &req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	// UPDATE: Mantiene mensaje
	utils.SendSuccess(c, http.StatusOK, "Usuario actualizado correctamente", gin.H{"user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de usuario inválido"))
		return
	}

	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	// DELETE: Mantiene mensaje
	utils.SendSuccess(c, http.StatusOK, "Usuario eliminado correctamente", nil)
}

// CheckUsernameAvailability maneja la verificación de username
func (h *UserHandler) CheckUsernameAvailability(c *gin.Context) {
	username := c.Query("username")
	
	// Validar que el username no esté vacío
	if username == "" {
		utils.HandleGinError(c, utils.NewBadRequestError("Username es requerido"))
		return
	}
	
	// Validar longitud mínima
	if len(username) < 3 {
		utils.HandleGinError(c, utils.NewBadRequestError("Username debe tener al menos 3 caracteres"))
		return
	}
	
	available, err := h.userService.CheckUsernameAvailability(username)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}
	
	// Respuesta simple que espera el frontend
	utils.SendData(c, http.StatusOK, gin.H{
		"available": available,
	})
}

// CheckEmailAvailability maneja la verificación de email
func (h *UserHandler) CheckEmailAvailability(c *gin.Context) {
	email := c.Query("email")
	
	// Validar que el email no esté vacío
	if email == "" {
		utils.HandleGinError(c, utils.NewBadRequestError("Email es requerido"))
		return
	}
	
	// Validación básica de formato email
	if !strings.Contains(email, "@") {
		utils.HandleGinError(c, utils.NewBadRequestError("Formato de email inválido"))
		return
	}
	
	available, err := h.userService.CheckEmailAvailability(email)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}
	
	// Respuesta simple que espera el frontend
	utils.SendData(c, http.StatusOK, gin.H{
		"available": available,
	})
}

