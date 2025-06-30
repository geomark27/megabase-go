package app_errors

import "fmt"

// ErrNotFound es una variable de error base para recursos no encontrados.
// Se define como 'var' para que podamos usar 'errors.Is' y comprobar si un error
// es de esta "familia".
var ErrNotFound = fmt.Errorf("resource not found")

// ErrConflict es un tipo de error estructurado para conflictos de datos (ej. valores duplicados).
// Se define como 'struct' para que podamos usar 'errors.As' y extraer los detalles.
type ErrConflict struct {
    Details string
}

// El método Error() hace que ErrConflict cumpla con la interfaz de 'error' estándar de Go.
func (e *ErrConflict) Error() string {
    return e.Details
}

// --- Funciones "Constructoras" ---
// Estas funciones ayudan a crear nuestros errores personalizados de una manera limpia.

// NewNotFoundError crea un error que "envuelve" a nuestro ErrNotFound base,
// añadiéndole más contexto sobre qué recurso no se encontró.
func NewNotFoundError(resource string, id uint) error {
    return fmt.Errorf("%w: %s with id %d not found", ErrNotFound, resource, id)
}

// NewConflictError crea una nueva instancia de nuestro error de conflicto estructurado.
func NewConflictError(details string) error {
    return &ErrConflict{Details: details}
}