package dto

import (
	"time"
	"gorm.io/datatypes"
)

// CreateCitizenRequest estructura para crear un citizen
// Piensa en esto como un formulario que debe llenarse para registrar un nuevo contribuyente
type CreateCitizenRequest struct {
	// --- IDENTIFICACIÓN PRINCIPAL (Obligatorio) ---
	NumeroIdentificacion string `json:"numero_identificacion" binding:"required,min=10,max=25"`
	TipoIdentificacion   string `json:"tipo_identificacion" binding:"required,oneof=04 05 06 07"`
	
	// --- DATOS DE CONTACTO (Obligatorio) ---
	Email              string `json:"email" binding:"required,email,max=100"`
	Celular            string `json:"celular" binding:"max=20"`
	Convencional       string `json:"convencional" binding:"max=20"`
	DireccionPrincipal string `json:"direccion_principal" binding:"max=250"`
	Pais               string `json:"pais" binding:"max=100"`
	Provincia          string `json:"provincia" binding:"max=100"`
	Ciudad             string `json:"ciudad" binding:"max=100"`

	// --- DATOS DE PERSONA NATURAL (Opcionales - usar cuando tipo_identificacion es 05 o 06) ---
	Nombre          *string    `json:"nombre,omitempty" binding:"omitempty,max=100"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	Nacionalidad    *string    `json:"nacionalidad,omitempty" binding:"omitempty,max=100"`
	EstadoCivil     *string    `json:"estado_civil,omitempty" binding:"omitempty,max=50"`
	Genero          *string    `json:"genero,omitempty" binding:"omitempty,max=50"`

	// --- DATOS DE EMPRESA (Opcionales - usar cuando tipo_identificacion es 04) ---
	RazonSocial           *string        `json:"razon_social,omitempty" binding:"omitempty,max=250"`
	NombreComercial       *string        `json:"nombre_comercial,omitempty" binding:"omitempty,max=250"`
	TipoEmpresa           *string        `json:"tipo_empresa,omitempty" binding:"omitempty,max=100"`
	RepresentantesLegales datatypes.JSON `json:"representantes_legales,omitempty"`

	// --- INFORMACIÓN TRIBUTARIA (Obligatorio) ---
	TipoContribuyente           string         `json:"tipo_contribuyente" binding:"required,max=100"`
	EstadoContribuyente         string         `json:"estado_contribuyente" binding:"required,oneof=ACTIVO SUSPENDIDO CANCELADO"`
	Regimen                     string         `json:"regimen" binding:"required,max=100"`
	Categoria                   string         `json:"categoria" binding:"max=100"`
	ObligadoContabilidad        string         `json:"obligado_contabilidad" binding:"required,oneof=SI NO"`
	AgenteRetencion             *string        `json:"agente_retencion,omitempty" binding:"omitempty,max=100"`
	ContribuyenteEspecial       *string        `json:"contribuyente_especial,omitempty" binding:"omitempty,max=100"`
	ActividadEconomicaPrincipal string         `json:"actividad_economica_principal" binding:"required,max=200"`
	Sucursales                  datatypes.JSON `json:"sucursales,omitempty"`

	// --- METADATOS ADICIONALES ---
	MotivoCancelacionSuspension string `json:"motivo_cancelacion_suspension,omitempty" binding:"max=250"`
}

// UpdateCitizenRequest estructura para actualizar un citizen
// Esta versión permite actualizaciones parciales - no todos los campos son obligatorios
type UpdateCitizenRequest struct {
	// --- IDENTIFICACIÓN PRINCIPAL ---
	NumeroIdentificacion *string `json:"numero_identificacion,omitempty" binding:"omitempty,min=10,max=25"`
	TipoIdentificacion   *string `json:"tipo_identificacion,omitempty" binding:"omitempty,oneof=04 05 06 07"`
	
	// --- DATOS DE CONTACTO ---
	Email              *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Celular            *string `json:"celular,omitempty" binding:"omitempty,max=20"`
	Convencional       *string `json:"convencional,omitempty" binding:"omitempty,max=20"`
	DireccionPrincipal *string `json:"direccion_principal,omitempty" binding:"omitempty,max=250"`
	Pais               *string `json:"pais,omitempty" binding:"omitempty,max=100"`
	Provincia          *string `json:"provincia,omitempty" binding:"omitempty,max=100"`
	Ciudad             *string `json:"ciudad,omitempty" binding:"omitempty,max=100"`

	// --- DATOS DE PERSONA NATURAL ---
	Nombre          *string    `json:"nombre,omitempty" binding:"omitempty,max=100"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	Nacionalidad    *string    `json:"nacionalidad,omitempty" binding:"omitempty,max=100"`
	EstadoCivil     *string    `json:"estado_civil,omitempty" binding:"omitempty,max=50"`
	Genero          *string    `json:"genero,omitempty" binding:"omitempty,max=50"`

	// --- DATOS DE EMPRESA ---
	RazonSocial           *string        `json:"razon_social,omitempty" binding:"omitempty,max=250"`
	NombreComercial       *string        `json:"nombre_comercial,omitempty" binding:"omitempty,max=250"`
	TipoEmpresa           *string        `json:"tipo_empresa,omitempty" binding:"omitempty,max=100"`
	RepresentantesLegales datatypes.JSON `json:"representantes_legales,omitempty"`

	// --- INFORMACIÓN TRIBUTARIA ---
	TipoContribuyente           *string        `json:"tipo_contribuyente,omitempty" binding:"omitempty,max=100"`
	EstadoContribuyente         *string        `json:"estado_contribuyente,omitempty" binding:"omitempty,oneof=ACTIVO SUSPENDIDO CANCELADO"`
	Regimen                     *string        `json:"regimen,omitempty" binding:"omitempty,max=100"`
	Categoria                   *string        `json:"categoria,omitempty" binding:"omitempty,max=100"`
	ObligadoContabilidad        *string        `json:"obligado_contabilidad,omitempty" binding:"omitempty,oneof=SI NO"`
	AgenteRetencion             *string        `json:"agente_retencion,omitempty" binding:"omitempty,max=100"`
	ContribuyenteEspecial       *string        `json:"contribuyente_especial,omitempty" binding:"omitempty,max=100"`
	ActividadEconomicaPrincipal *string        `json:"actividad_economica_principal,omitempty" binding:"omitempty,max=200"`
	Sucursales                  datatypes.JSON `json:"sucursales,omitempty"`

	// --- METADATOS ADICIONALES ---
	MotivoCancelacionSuspension *string `json:"motivo_cancelacion_suspension,omitempty" binding:"omitempty,max=250"`
}

// CitizenResponse estructura para respuestas
// Esta es la información que se devuelve al cliente, incluyendo campos calculados
type CitizenResponse struct {
	ID uint `json:"id"`

	// --- IDENTIFICACIÓN PRINCIPAL ---
	NumeroIdentificacion string `json:"numero_identificacion"`
	TipoIdentificacion   string `json:"tipo_identificacion"`
	
	// --- DATOS DE CONTACTO ---
	Email              string `json:"email"`
	Celular            string `json:"celular"`
	Convencional       string `json:"convencional"`
	DireccionPrincipal string `json:"direccion_principal"`
	Pais               string `json:"pais"`
	Provincia          string `json:"provincia"`
	Ciudad             string `json:"ciudad"`

	// --- DATOS DE PERSONA NATURAL ---
	Nombre          *string    `json:"nombre,omitempty"`
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	Nacionalidad    *string    `json:"nacionalidad,omitempty"`
	EstadoCivil     *string    `json:"estado_civil,omitempty"`
	Genero          *string    `json:"genero,omitempty"`
	
	// Campo calculado: edad basada en fecha de nacimiento
	Edad *int `json:"edad,omitempty"`

	// --- DATOS DE EMPRESA ---
	RazonSocial           *string        `json:"razon_social,omitempty"`
	NombreComercial       *string        `json:"nombre_comercial,omitempty"`
	TipoEmpresa           *string        `json:"tipo_empresa,omitempty"`
	RepresentantesLegales datatypes.JSON `json:"representantes_legales,omitempty"`

	// --- INFORMACIÓN TRIBUTARIA ---
	TipoContribuyente           string         `json:"tipo_contribuyente"`
	EstadoContribuyente         string         `json:"estado_contribuyente"`
	Regimen                     string         `json:"regimen"`
	Categoria                   string         `json:"categoria"`
	ObligadoContabilidad        string         `json:"obligado_contabilidad"`
	AgenteRetencion             *string        `json:"agente_retencion,omitempty"`
	ContribuyenteEspecial       *string        `json:"contribuyente_especial,omitempty"`
	ActividadEconomicaPrincipal string         `json:"actividad_economica_principal"`
	Sucursales                  datatypes.JSON `json:"sucursales,omitempty"`

	// --- METADATOS ---
	MotivoCancelacionSuspension string      `json:"motivo_cancelacion_suspension,omitempty"`
	CreatedAt                   interface{} `json:"created_at"`
	UpdatedAt                   interface{} `json:"updated_at"`
}

// CitizenSearchFilters estructura para filtros de búsqueda
// Esto permite hacer búsquedas más específicas y complejas
type CitizenSearchFilters struct {
	TipoIdentificacion  *string `form:"tipo_identificacion" binding:"omitempty,oneof=04 05 06 07"`
	EstadoContribuyente *string `form:"estado_contribuyente" binding:"omitempty,oneof=ACTIVO SUSPENDIDO CANCELADO"`
	Regimen             *string `form:"regimen"`
	Pais                *string `form:"pais"`
	Provincia           *string `form:"provincia"`
	Ciudad              *string `form:"ciudad"`
	ObligadoContabilidad *string `form:"obligado_contabilidad" binding:"omitempty,oneof=SI NO"`
	
	// Paginación
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=10" binding:"min=1,max=100"`
}