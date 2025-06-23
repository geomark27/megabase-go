package handlers

import (
	"net/http"
	"strconv"

	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"
	"megabaseGo/internal/utils"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *services.RoleService
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: services.NewRoleService(),
	}
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusCreated, "Rol creado exitosamente", gin.H{"role": role})
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	includeInactive := c.Query("include_inactive") == "true"

	roles, err := h.roleService.GetRoles(includeInactive)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Roles obtenidos correctamente", gin.H{
		"roles": roles,
		"count": len(roles),
	})
}

func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de rol inválido"))
		return
	}

	role, err := h.roleService.GetRoleByID(uint(roleID))
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Rol obtenido correctamente", gin.H{"role": role})
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de rol inválido"))
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleGinError(c, err)
		return
	}

	role, err := h.roleService.UpdateRole(uint(roleID), &req)
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Rol actualizado correctamente", gin.H{"role": role})
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	roleID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.HandleGinError(c, utils.NewBadRequestError("ID de rol inválido"))
		return
	}

	err = h.roleService.DeleteRole(uint(roleID))
	if err != nil {
		utils.HandleGinError(c, err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Rol eliminado correctamente", nil)
}