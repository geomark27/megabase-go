package handlers

import (
	"errors"
	"megabaseGo/internal/app/dto"
	app_errors "megabaseGo/internal/app/errors"
	"megabaseGo/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	svc *services.CompanyService
}

func NewCompanyHandler() *CompanyHandler {
	return &CompanyHandler{
		svc: services.NewCompanyService(),
	}
}

func (h *CompanyHandler) handleError(c *gin.Context, err error) {
	var conflictErr *app_errors.ErrConflict

	if errors.As(err, &conflictErr) {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Data conflict",
			"details": conflictErr.Error(),
		})
		return
	}

	if errors.Is(err, app_errors.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Resource not found",
			"details": err.Error(),
		})
		return
	}

	// Si no es ninguno de los anteriores, es un error genérico del servidor.
	c.JSON(http.StatusInternalServerError, gin.H{
		"error":   "Internal server error",
		"details": err.Error(),
	})
}

// --- El resto de los handlers usan la nueva función handleError ---

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req dto.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	company, err := h.svc.CreateCompany(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": company})
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	company, err := h.svc.GetCompanyByID(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": company})
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req dto.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}
	company, err := h.svc.UpdateCompany(uint(id), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": company})
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	err = h.svc.DeleteCompany(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Company deleted successfully"})
}

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	var filters dto.CompanySearchFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}
	companies, err := h.svc.GetCompanies(&filters)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    companies,
		"count":   len(companies),
		"filters": filters,
	})
}