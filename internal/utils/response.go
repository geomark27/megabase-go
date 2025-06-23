package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonResponse define la estructura estándar para todas las respuestas de la API.
type JsonResponse struct {
	Status  string      `json:"status"`            // "success" o "error"
	Message string      `json:"message"`           // Mensaje descriptivo
	Data    interface{} `json:"data,omitempty"`    // Datos de la respuesta (opcional)
}

// SendSuccess es el nuevo helper para enviar respuestas exitosas.
// Sigue el patrón de Laravel: código, mensaje y datos.
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := JsonResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// SendError es el nuevo helper para enviar respuestas de error consistentes.
func SendError(c *gin.Context, statusCode int, message string) {
	response := JsonResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
	c.JSON(statusCode, response)
}

// HandleGinError maneja los errores de validación de Gin y otros errores genéricos.
// Usaremos este para reemplazar a HandleValidationError y HandleError.
func HandleGinError(c *gin.Context, err error) {
	// Intentamos ver si es un error de nuestra API que ya hemos definido
	if apiErr, ok := IsAPIError(err); ok {
		SendError(c, apiErr.GetStatusCode(), apiErr.Error())
		return
	}

	// Si no, es probablemente un error de validación de Gin o un error 500
	// Aquí podrías añadir un logging más detallado si quisieras
	SendError(c, http.StatusBadRequest, err.Error())
}