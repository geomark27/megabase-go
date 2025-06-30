package handlers

import (
    "megabaseGo/internal/app/dto"
    "megabaseGo/internal/app/services"
    "net/http"

    "github.com/gin-gonic/gin"
)

type ConsultHandler struct {
    svc *services.ConsultService
}

func NewConsultHandler() *ConsultHandler {
    return &ConsultHandler{
        svc: services.NewConsultService(),
    }
}

// ConsultHandler maneja las consultas de identificación
// @Tags Consult
// @Summary Consulta información de cédula o RUC
// @Description Valida el número de identificación y realiza la consulta externa según tipo (cédula o RUC)
// @Accept json
// @Produce json
// @Param request body dto.ConsultRequest true "Datos de consulta"   
// @Success 200 {object} interface{} "JSON devuelto por la API externa o respuesta de validación"
// @Failure 400 {object} map[string]string "Invalid request data"
// @Failure 500 {object} map[string]string "Consulta fallida"
// @Router /api/v1/consult [post]
func (h *ConsultHandler) Consultar(c *gin.Context) {
	var req dto.ConsultRequest
	// Bind del JSON al struct, con validación de campos
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Llamada al servicio que maneja validaciones y consumo externo
	resp, err := h.svc.GetCitizenByNumeroIdentificacion(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Consulta fallida", "details": err.Error()})
		return
	}

	// Respuesta exitosa: se devuelve directamente el JSON resultante
	c.JSON(http.StatusOK, resp)
}
