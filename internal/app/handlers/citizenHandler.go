package handlers

import (
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CitizenHandler maneja todas las peticiones HTTP relacionadas con ciudadanos
type CitizenHandler struct {
	citizenService *services.CitizenService
}

// NewCitizenHandler crea una nueva instancia del handler
func NewCitizenHandler() *CitizenHandler {
	return &CitizenHandler{
		citizenService: services.NewCitizenService(),
	}
}

// handleError es una función helper para centralizar el manejo de errores de ciudadano.
// Determina el código de estado HTTP correcto y formatea la respuesta JSON.
func (h *CitizenHandler) handleError(c *gin.Context, err error, defaultMessage string, defaultStatus int) {
	statusCode := defaultStatus
	errorMessage := defaultMessage

	errStr := err.Error()

	if strings.Contains(errStr, "not found") {
		statusCode = http.StatusNotFound
		errorMessage = "Resource not found"
	} else if strings.Contains(errStr, "already exists") ||
		strings.Contains(errStr, "duplicate") ||
		strings.Contains(errStr, "ya esta registrado") ||
		strings.Contains(errStr, "requires") {
		statusCode = http.StatusConflict // <-- El código correcto para conflictos de datos.
		errorMessage = "Data validation error or conflict"
	}

	c.JSON(statusCode, gin.H{
		"error":   errorMessage,
		"details": errStr, // Siempre enviamos el error original y específico en 'details'.
	})
}

// GetAllCitizens maneja GET /citizens con filtros opcionales
func (h *CitizenHandler) GetAllCitizens(c *gin.Context) {
	var filters dto.CitizenSearchFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	citizens, err := h.citizenService.GetAllCitizens(&filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve citizens",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizens,
		"count":   len(citizens),
		"filters": filters,
	})
}

// GetCitizenByID maneja GET /citizens/:id
func (h *CitizenHandler) GetCitizenByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid citizen ID",
			"details": "ID must be a positive number",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByID(uint(id))
	if err != nil {
		h.handleError(c, err, "Failed to retrieve citizen", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByEmail maneja GET /citizens/email/:email
func (h *CitizenHandler) GetCitizenByEmail(c *gin.Context) {
	email := c.Param("email")

	if email == "" || !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid email format",
			"details": "Please provide a valid email address",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByEmail(email)
	if err != nil {
		h.handleError(c, err, "Failed to retrieve citizen by email", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByIdentification maneja GET /citizens/identification/:numero
func (h *CitizenHandler) GetCitizenByIdentification(c *gin.Context) {
	numero := c.Param("numero")

	if numero == "" || len(numero) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid identification number",
			"details": "Identification number must be at least 10 characters",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByNumeroIdentificacion(numero)
	if err != nil {
		h.handleError(c, err, "Failed to retrieve citizen by identification", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByRazonSocial maneja GET /citizens/razon-social/:razon
func (h *CitizenHandler) GetCitizenByRazonSocial(c *gin.Context) {
	razonSocial := c.Param("razon")

	if razonSocial == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid razon social",
			"details": "Razon social cannot be empty",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByRazonSocial(razonSocial)
	if err != nil {
		h.handleError(c, err, "Failed to retrieve citizen by razon social", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// CreateCitizen maneja POST /citizens
func (h *CitizenHandler) CreateCitizen(c *gin.Context) {
	var req dto.CreateCitizenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	citizen, err := h.citizenService.CreateCitizen(&req)
	if err != nil {
		h.handleError(c, err, "Failed to create citizen", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Citizen created successfully",
		"data":    citizen,
	})
}

// UpdateCitizen maneja PUT /citizens/:id
func (h *CitizenHandler) UpdateCitizen(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid citizen ID",
			"details": "ID must be a positive number",
		})
		return
	}

	var req dto.UpdateCitizenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	citizen, err := h.citizenService.UpdateCitizen(uint(id), &req)
	if err != nil {
		h.handleError(c, err, "Failed to update citizen", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Citizen updated successfully",
		"data":    citizen,
	})
}

// DeleteCitizen maneja DELETE /citizens/:id
func (h *CitizenHandler) DeleteCitizen(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid citizen ID",
			"details": "ID must be a positive number",
		})
		return
	}

	err = h.citizenService.DeleteCitizen(uint(id))
	if err != nil {
		h.handleError(c, err, "Failed to delete citizen", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Citizen deleted successfully",
	})
}

// --- ENDPOINTS DE VERIFICACIÓN ---

// CheckIdentificationAvailability maneja GET /citizens/check/identification/:numero
func (h *CitizenHandler) CheckIdentificationAvailability(c *gin.Context) {
	numero := c.Param("numero")

	if numero == "" || len(numero) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid identification number",
		})
		return
	}

	_, err := h.citizenService.GetCitizenByNumeroIdentificacion(numero)
	available := err != nil && strings.Contains(err.Error(), "not found")

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"numero":    numero,
	})
}

// CheckEmailAvailability maneja GET /citizens/check/email/:email
func (h *CitizenHandler) CheckEmailAvailability(c *gin.Context) {
	email := c.Param("email")

	if email == "" || !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email format",
		})
		return
	}

	_, err := h.citizenService.GetCitizenByEmail(email)
	available := err != nil && strings.Contains(err.Error(), "not found")

	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"email":     email,
	})
}

// CheckRazonSocialAvailability maneja GET /citizens/check/razon-social/:razon
func (h *CitizenHandler) CheckRazonSocialAvailability(c *gin.Context) {
	razonSocial := c.Param("razon")

	if razonSocial == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid razon social",
		})
		return
	}

	_, err := h.citizenService.GetCitizenByRazonSocial(razonSocial)
	available := err != nil && strings.Contains(err.Error(), "not found")

	c.JSON(http.StatusOK, gin.H{
		"available":    available,
		"razon_social": razonSocial,
	})
}