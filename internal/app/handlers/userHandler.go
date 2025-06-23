package handlers

import (
	"net/http"
	"strconv"

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
		// Usamos el nuevo manejador de errores
		utils.HandleGinError(c, err)
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		// También lo usamos para errores del servicio
		utils.HandleGinError(c, err)
		return
	}

	// Usamos el nuevo helper de éxito, ¡igual que en Laravel!
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

	// Usamos SendSuccess también aquí
	utils.SendSuccess(c, http.StatusOK, "Usuarios obtenidos correctamente", gin.H{
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

	utils.SendSuccess(c, http.StatusOK, "Usuario obtenido correctamente", gin.H{"user": user})
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

	// Para respuestas sin datos, simplemente pasamos nil
	utils.SendSuccess(c, http.StatusOK, "Usuario eliminado correctamente", nil)
}