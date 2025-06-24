package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Contribuyente representa tanto a una Persona Natural como a una Sociedad (Persona Jurídica)
// registrada en el sistema, principalmente para fines fiscales y de facturación.
type Citizen struct {
	// --- CAMPOS DE AUDITORÍA (GORM) ---
	gorm.Model // Incluye ID, CreatedAt, UpdatedAt, DeletedAt

	// --- 1. IDENTIFICACIÓN PRINCIPAL (Ambos tipos) ---
	// Es el campo más importante. Se recomienda un tamaño más ajustado.
	// 13 para RUC, 10 para Cédula. 25 es un tamaño seguro para incluir pasaportes, etc.
	NumeroIdentificacion string `gorm:"size:25;not null;uniqueIndex" json:"numero_identificacion"`

	// Códigos del SRI: '04' (RUC), '05' (Cédula), '06' (Pasaporte)
	TipoIdentificacion string `gorm:"size:2;not null;check:tipo_identificacion IN ('04','05','06','07')" json:"tipo_identificacion"`
	
	// --- 2. DATOS DE CONTACTO Y UBICACIÓN (Ambos tipos) ---
	Email        string `gorm:"size:100;not null" json:"email"`
	Celular      string `gorm:"size:20" json:"celular"`
	Convencional string `gorm:"size:20" json:"convencional"`
	// Dirección principal o domicilio fiscal
	DireccionPrincipal string `gorm:"size:250" json:"direccion_principal"`
	Pais               string `gorm:"size:100" json:"pais"`
	Provincia          string `gorm:"size:100" json:"provincia"`
	Ciudad             string `gorm:"size:100" json:"ciudad"`


	// --- 3. DATOS EXCLUSIVOS DE PERSONA NATURAL ---
	// Estos campos deben ser punteros (*) para permitir valores NULOS si el contribuyente es una empresa.
	Nombre          *string    `gorm:"size:100" json:"nombre,omitempty"` // Nombre completo (Apellido y Nombre)
	FechaNacimiento *time.Time `json:"fecha_nacimiento,omitempty"`
	Nacionalidad    *string    `gorm:"size:100" json:"nacionalidad,omitempty"`
	EstadoCivil     *string    `gorm:"size:50" json:"estado_civil,omitempty"`
	Genero          *string    `gorm:"size:50" json:"genero,omitempty"`
	// El campo "Edad" se ha eliminado. Es un dato calculado (Ahora - FechaNacimiento) y no debe almacenarse.

	// --- 4. DATOS EXCLUSIVOS DE SOCIEDAD / EMPRESA ---
	// Estos campos deben ser punteros (*) para permitir valores NULOS si el contribuyente es una persona natural.
	RazonSocial          *string `gorm:"size:250;uniqueIndex" json:"razon_social,omitempty"`
	NombreComercial      *string `gorm:"size:250;uniqueIndex" json:"nombre_comercial,omitempty"`
	TipoEmpresa          *string `gorm:"size:100" json:"tipo_empresa,omitempty"`
	RepresentantesLegales datatypes.JSON `json:"representantes_legales,omitempty"`


	// --- 5. INFORMACIÓN TRIBUTARIA (SRI - Ambos tipos) ---
	// La mayoría de estos datos provienen de la consulta al SRI
	TipoContribuyente           string         `gorm:"size:100" json:"tipo_contribuyente"`
	EstadoContribuyente         string         `gorm:"size:100" json:"estado_contribuyente"` // ACTIVO, SUSPENDIDO, CANCELADO
	Regimen                     string         `gorm:"size:100" json:"regimen"` // Ej: Régimen General, RIMPE
	Categoria                   string         `gorm:"size:100" json:"categoria,omitempty"`
	ObligadoContabilidad        string         `gorm:"size:2" json:"obligado_contabilidad"` // SI / NO
	AgenteRetencion             *string        `gorm:"size:100" json:"agente_retencion,omitempty"`
	ContribuyenteEspecial       *string        `gorm:"size:100" json:"contribuyente_especial,omitempty"`
	ActividadEconomicaPrincipal string         `gorm:"size:200" json:"actividad_economica_principal"`
	Sucursales                  datatypes.JSON `json:"sucursales,omitempty"`

	// --- 6. METADATOS ADICIONALES ---
	// Corregido typo: "suspencion" a "suspension"
	MotivoCancelacionSuspension string `gorm:"size:250" json:"motivo_cancelacion_suspension,omitempty"`
}