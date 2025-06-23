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

// JsonDataResponse define la estructura simple para respuestas GET (solo status y data)
type JsonDataResponse struct {
	Status string      `json:"status"` // "success" o "error"
	Data   interface{} `json:"data"`   // Datos de la respuesta
}

// SendSuccess es el helper para enviar respuestas exitosas con mensaje.
// Se usa para CREATE, UPDATE, DELETE
func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := JsonResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// SendData es el helper para enviar respuestas GET simples (solo status y data).
// Se usa para métodos GET que obtienen datos
func SendData(c *gin.Context, statusCode int, data interface{}) {
	response := JsonDataResponse{
		Status: "success",
		Data:   data,
	}
	c.JSON(statusCode, response)
}

// SendError es el helper para enviar respuestas de error consistentes.
func SendError(c *gin.Context, statusCode int, message string) {
	response := JsonResponse{
		Status:  "error",
		Message: message,
		Data:    nil,
	}
	c.JSON(statusCode, response)
}

// HandleGinError maneja los errores de validación de Gin y otros errores genéricos.
func HandleGinError(c *gin.Context, err error) {
	// Intentamos ver si es un error de nuestra API que ya hemos definido
	if apiErr, ok := IsAPIError(err); ok {
		SendError(c, apiErr.GetStatusCode(), apiErr.Error())
		return
	}

	// Si no, es probablemente un error de validación de Gin o un error 500
	SendError(c, http.StatusBadRequest, err.Error())
}