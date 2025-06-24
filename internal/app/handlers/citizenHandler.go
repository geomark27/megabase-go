package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/app/services"
	
	"github.com/gin-gonic/gin"
)

// CitizenHandler maneja todas las peticiones HTTP relacionadas con ciudadanos
// Este handler actúa como el "traductor" entre el mundo HTTP y la lógica de negocio
type CitizenHandler struct {
	citizenService *services.CitizenService
}

// NewCitizenHandler crea una nueva instancia del handler
func NewCitizenHandler() *CitizenHandler {
	return &CitizenHandler{
		citizenService: services.NewCitizenService(),
	}
}

// GetAllCitizens maneja GET /citizens con filtros opcionales
// @Summary Lista todos los ciudadanos con filtros
// @Description Obtiene una lista paginada de ciudadanos con filtros opcionales
// @Tags Citizens
// @Accept json
// @Produce json
// @Param tipo_identificacion query string false "Tipo de identificación (04,05,06,07)"
// @Param estado_contribuyente query string false "Estado del contribuyente (ACTIVO,SUSPENDIDO,CANCELADO)"
// @Param regimen query string false "Régimen tributario"
// @Param pais query string false "País"
// @Param provincia query string false "Provincia"
// @Param ciudad query string false "Ciudad"
// @Param obligado_contabilidad query string false "Obligado contabilidad (SI,NO)"
// @Param page query int false "Número de página" default(1)
// @Param page_size query int false "Tamaño de página" default(10)
// @Success 200 {object} map[string]interface{} "Lista de ciudadanos"
// @Failure 400 {object} map[string]interface{} "Error en parámetros"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens [get]
func (h *CitizenHandler) GetAllCitizens(c *gin.Context) {
	// Parsear filtros desde query parameters
	// Gin automáticamente convierte los query params a nuestra estructura
	var filters dto.CitizenSearchFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Llamar al service para obtener los datos
	citizens, err := h.citizenService.GetAllCitizens(&filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve citizens",
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa con metadatos útiles
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizens,
		"count":   len(citizens),
		"filters": filters,
	})
}

// GetCitizenByID maneja GET /citizens/:id
// @Summary Obtiene un ciudadano por ID
// @Description Obtiene los detalles completos de un ciudadano específico
// @Tags Citizens
// @Accept json
// @Produce json
// @Param id path int true "ID del ciudadano"
// @Success 200 {object} dto.CitizenResponse "Detalles del ciudadano"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 404 {object} map[string]interface{} "Ciudadano no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/{id} [get]
func (h *CitizenHandler) GetCitizenByID(c *gin.Context) {
	// Extraer y validar el ID del path parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid citizen ID",
			"details": "ID must be a positive number",
		})
		return
	}

	// Buscar el ciudadano
	citizen, err := h.citizenService.GetCitizenByID(uint(id))
	if err != nil {
		// Distinguir entre "no encontrado" y "error del servidor"
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Citizen not found",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve citizen",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByEmail maneja GET /citizens/email/:email
// @Summary Busca un ciudadano por email
// @Description Busca un ciudadano específico usando su dirección de email
// @Tags Citizens
// @Accept json
// @Produce json
// @Param email path string true "Email del ciudadano"
// @Success 200 {object} dto.CitizenResponse "Detalles del ciudadano"
// @Failure 404 {object} map[string]interface{} "Ciudadano no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/email/{email} [get]
func (h *CitizenHandler) GetCitizenByEmail(c *gin.Context) {
	email := c.Param("email")
	
	// Validación básica del email
	if email == "" || !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid email format",
			"details": "Please provide a valid email address",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByEmail(email)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Citizen not found with this email",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve citizen",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByIdentification maneja GET /citizens/identification/:numero
// @Summary Busca un ciudadano por número de identificación
// @Description Busca un ciudadano usando su número de identificación fiscal
// @Tags Citizens
// @Accept json
// @Produce json
// @Param numero path string true "Número de identificación"
// @Success 200 {object} dto.CitizenResponse "Detalles del ciudadano"
// @Failure 400 {object} map[string]interface{} "Número inválido"
// @Failure 404 {object} map[string]interface{} "Ciudadano no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/identification/{numero} [get]
func (h *CitizenHandler) GetCitizenByIdentification(c *gin.Context) {
	numero := c.Param("numero")
	
	// Validación básica del número de identificación
	if numero == "" || len(numero) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid identification number",
			"details": "Identification number must be at least 10 characters",
		})
		return
	}

	citizen, err := h.citizenService.GetCitizenByNumeroIdentificacion(numero)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Citizen not found with this identification number",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve citizen",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// GetCitizenByRazonSocial maneja GET /citizens/razon-social/:razon
// @Summary Busca una empresa por razón social
// @Description Busca una empresa específica usando su razón social
// @Tags Citizens
// @Accept json
// @Produce json
// @Param razon path string true "Razón social de la empresa"
// @Success 200 {object} dto.CitizenResponse "Detalles de la empresa"
// @Failure 400 {object} map[string]interface{} "Razón social inválida"
// @Failure 404 {object} map[string]interface{} "Empresa no encontrada"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/razon-social/{razon} [get]
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
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Citizen not found with this razon social",
				"details": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to retrieve citizen",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    citizen,
	})
}

// CreateCitizen maneja POST /citizens
// @Summary Crea un nuevo ciudadano
// @Description Registra un nuevo ciudadano/contribuyente en el sistema
// @Tags Citizens
// @Accept json
// @Produce json
// @Param citizen body dto.CreateCitizenRequest true "Datos del ciudadano"
// @Success 201 {object} dto.CitizenResponse "Ciudadano creado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 409 {object} map[string]interface{} "Conflicto - datos duplicados"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens [post]
func (h *CitizenHandler) CreateCitizen(c *gin.Context) {
	var req dto.CreateCitizenRequest
	
	// Validar y parsear el JSON del request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Crear el ciudadano usando el service
	citizen, err := h.citizenService.CreateCitizen(&req)
	if err != nil {
		// Determinar el tipo de error y responder apropiadamente
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to create citizen"

		// Errores de validación/duplicación deberían ser 409 (Conflict)
		if strings.Contains(err.Error(), "already exists") ||
		   strings.Contains(err.Error(), "duplicate") ||
		   strings.Contains(err.Error(), "requires") {
			statusCode = http.StatusConflict
			errorMessage = "Data validation error"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorMessage,
			"details": err.Error(),
		})
		return
	}

	// Respuesta exitosa con código 201 (Created)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Citizen created successfully",
		"data":    citizen,
	})
}

// UpdateCitizen maneja PUT /citizens/:id
// @Summary Actualiza un ciudadano existente
// @Description Actualiza los datos de un ciudadano específico
// @Tags Citizens
// @Accept json
// @Produce json
// @Param id path int true "ID del ciudadano"
// @Param citizen body dto.UpdateCitizenRequest true "Datos a actualizar"
// @Success 200 {object} dto.CitizenResponse "Ciudadano actualizado exitosamente"
// @Failure 400 {object} map[string]interface{} "Datos inválidos"
// @Failure 404 {object} map[string]interface{} "Ciudadano no encontrado"
// @Failure 409 {object} map[string]interface{} "Conflicto - datos duplicados"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/{id} [put]
func (h *CitizenHandler) UpdateCitizen(c *gin.Context) {
	// Validar ID
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
	
	// Validar y parsear el JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Actualizar usando el service
	citizen, err := h.citizenService.UpdateCitizen(uint(id), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to update citizen"

		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
			errorMessage = "Citizen not found"
		} else if strings.Contains(err.Error(), "already exists") ||
				 strings.Contains(err.Error(), "duplicate") {
			statusCode = http.StatusConflict
			errorMessage = "Data validation error"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorMessage,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Citizen updated successfully",
		"data":    citizen,
	})
}

// DeleteCitizen maneja DELETE /citizens/:id
// @Summary Elimina un ciudadano
// @Description Realiza un soft delete de un ciudadano específico
// @Tags Citizens
// @Accept json
// @Produce json
// @Param id path int true "ID del ciudadano"
// @Success 200 {object} map[string]interface{} "Ciudadano eliminado exitosamente"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 404 {object} map[string]interface{} "Ciudadano no encontrado"
// @Failure 500 {object} map[string]interface{} "Error interno"
// @Router /citizens/{id} [delete]
func (h *CitizenHandler) DeleteCitizen(c *gin.Context) {
	// Validar ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid citizen ID",
			"details": "ID must be a positive number",
		})
		return
	}

	// Eliminar usando el service
	err = h.citizenService.DeleteCitizen(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to delete citizen"

		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
			errorMessage = "Citizen not found"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorMessage,
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Citizen deleted successfully",
	})
}

// --- ENDPOINTS DE VERIFICACIÓN ---
// Estos endpoints son útiles para validación en tiempo real en el frontend

// CheckIdentificationAvailability maneja GET /citizens/check/identification/:numero
// @Summary Verifica disponibilidad de número de identificación
// @Description Verifica si un número de identificación está disponible para uso
// @Tags Citizens
// @Accept json
// @Produce json
// @Param numero path string true "Número de identificación a verificar"
// @Success 200 {object} map[string]interface{} "Estado de disponibilidad"
// @Failure 400 {object} map[string]interface{} "Número inválido"
// @Router /citizens/check/identification/{numero} [get]
func (h *CitizenHandler) CheckIdentificationAvailability(c *gin.Context) {
	numero := c.Param("numero")
	
	if numero == "" || len(numero) < 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid identification number",
			"details": "Identification number must be at least 10 characters",
		})
		return
	}

	// Intentar buscar el ciudadano
	_, err := h.citizenService.GetCitizenByNumeroIdentificacion(numero)
	
	// Si encontramos algo, no está disponible
	available := err != nil && strings.Contains(err.Error(), "not found")
	
	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"numero":    numero,
	})
}

// CheckEmailAvailability maneja GET /citizens/check/email/:email
// @Summary Verifica disponibilidad de email
// @Description Verifica si un email está disponible para uso
// @Tags Citizens
// @Accept json
// @Produce json
// @Param email path string true "Email a verificar"
// @Success 200 {object} map[string]interface{} "Estado de disponibilidad"
// @Failure 400 {object} map[string]interface{} "Email inválido"
// @Router /citizens/check/email/{email} [get]
func (h *CitizenHandler) CheckEmailAvailability(c *gin.Context) {
	email := c.Param("email")
	
	if email == "" || !strings.Contains(email, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid email format",
			"details": "Please provide a valid email address",
		})
		return
	}

	// Intentar buscar el ciudadano
	_, err := h.citizenService.GetCitizenByEmail(email)
	
	// Si encontramos algo, no está disponible
	available := err != nil && strings.Contains(err.Error(), "not found")
	
	c.JSON(http.StatusOK, gin.H{
		"available": available,
		"email":     email,
	})
}

// CheckRazonSocialAvailability maneja GET /citizens/check/razon-social/:razon
// @Summary Verifica disponibilidad de razón social
// @Description Verifica si una razón social está disponible para uso
// @Tags Citizens
// @Accept json
// @Produce json
// @Param razon path string true "Razón social a verificar"
// @Success 200 {object} map[string]interface{} "Estado de disponibilidad"
// @Failure 400 {object} map[string]interface{} "Razón social inválida"
// @Router /citizens/check/razon-social/{razon} [get]
func (h *CitizenHandler) CheckRazonSocialAvailability(c *gin.Context) {
	razonSocial := c.Param("razon")
	
	if razonSocial == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid razon social",
			"details": "Razon social cannot be empty",
		})
		return
	}

	// Intentar buscar el ciudadano
	_, err := h.citizenService.GetCitizenByRazonSocial(razonSocial)
	
	// Si encontramos algo, no está disponible
	available := err != nil && strings.Contains(err.Error(), "not found")
	
	c.JSON(http.StatusOK, gin.H{
		"available":    available,
		"razon_social": razonSocial,
	})
}