package services

import (
	"errors"
	"time"
	"megabaseGo/internal/app/dto"
	"megabaseGo/internal/database"
	"megabaseGo/internal/models"
	"gorm.io/gorm"
)

// CitizenService maneja la lógica de negocio para ciudadanos/contribuyentes
// Piensa en este service como el "cerebro" que toma decisiones inteligentes
// sobre cómo manejar los datos de contribuyentes ecuatorianos
type CitizenService struct{}

// NewCitizenService crea una nueva instancia del servicio
func NewCitizenService() *CitizenService {
	return &CitizenService{}
}

// GetAllCitizens obtiene todos los ciudadanos con filtros opcionales
// Este método es inteligente - permite filtrar por múltiples criterios
func (s *CitizenService) GetAllCitizens(filters *dto.CitizenSearchFilters) ([]dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizens []models.Citizen

	// Construir la query base
	query := db.Model(&models.Citizen{})

	// Aplicar filtros si están presentes
	// Esto es como construir una búsqueda personalizada paso a paso
	if filters != nil {
		if filters.TipoIdentificacion != nil {
			query = query.Where("tipo_identificacion = ?", *filters.TipoIdentificacion)
		}
		if filters.EstadoContribuyente != nil {
			query = query.Where("estado_contribuyente = ?", *filters.EstadoContribuyente)
		}
		if filters.Regimen != nil {
			query = query.Where("regimen ILIKE ?", "%"+*filters.Regimen+"%")
		}
		if filters.Pais != nil {
			query = query.Where("pais ILIKE ?", "%"+*filters.Pais+"%")
		}
		if filters.Provincia != nil {
			query = query.Where("provincia ILIKE ?", "%"+*filters.Provincia+"%")
		}
		if filters.Ciudad != nil {
			query = query.Where("ciudad ILIKE ?", "%"+*filters.Ciudad+"%")
		}
		if filters.ObligadoContabilidad != nil {
			query = query.Where("obligado_contabilidad = ?", *filters.ObligadoContabilidad)
		}

		// Aplicar paginación
		// Esto es importante para no sobrecargar el sistema con muchos resultados
		offset := (filters.Page - 1) * filters.PageSize
		query = query.Offset(offset).Limit(filters.PageSize)
	}

	// Ejecutar la consulta
	if err := query.Find(&citizens).Error; err != nil {
		return nil, err
	}

	// Convertir a DTOs de respuesta
	var responses []dto.CitizenResponse
	for _, citizen := range citizens {
		responses = append(responses, *s.toCitizenResponse(&citizen))
	}

	return responses, nil
}

// GetCitizenByID obtiene un ciudadano por su ID
func (s *CitizenService) GetCitizenByID(id uint) (*dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizen models.Citizen

	if err := db.First(&citizen, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("citizen not found")
		}
		return nil, err
	}

	return s.toCitizenResponse(&citizen), nil
}

// GetCitizenByEmail busca un ciudadano por email
// Email debe ser único en el sistema
func (s *CitizenService) GetCitizenByEmail(email string) (*dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizen models.Citizen

	if err := db.Where("email = ?", email).First(&citizen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("citizen not found with this email")
		}
		return nil, err
	}

	return s.toCitizenResponse(&citizen), nil
}

// GetCitizenByNumeroIdentificacion busca por número de identificación
// Este es el método más importante - la identificación fiscal es única por ley
func (s *CitizenService) GetCitizenByNumeroIdentificacion(numero string) (*dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizen models.Citizen

	if err := db.Where("numero_identificacion = ?", numero).First(&citizen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("citizen not found with this identification number")
		}
		return nil, err
	}

	return s.toCitizenResponse(&citizen), nil
}

// GetCitizenByRazonSocial busca empresas por razón social
// Solo aplica para empresas (RUC), no para personas naturales
func (s *CitizenService) GetCitizenByRazonSocial(razonSocial string) (*dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizen models.Citizen

	if err := db.Where("razon_social = ?", razonSocial).First(&citizen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("citizen not found with this razon social")
		}
		return nil, err
	}

	return s.toCitizenResponse(&citizen), nil
}

// CreateCitizen crea un nuevo ciudadano con todas las validaciones necesarias
// Este método es el corazón del sistema - debe validar todo cuidadosamente
func (s *CitizenService) CreateCitizen(req *dto.CreateCitizenRequest) (*dto.CitizenResponse, error) {
	db := database.GetDB()

	// VALIDACIÓN 1: Verificar que el número de identificación no exista
	// Esta es la validación más crítica - no puede haber duplicados fiscales
	if err := s.validateUniqueNumeroIdentificacion(req.NumeroIdentificacion, 0); err != nil {
		return nil, err
	}

	// VALIDACIÓN 2: Verificar que el email no exista
	if err := s.validateUniqueEmail(req.Email, 0); err != nil {
		return nil, err
	}

	// VALIDACIÓN 3: Si es empresa, verificar que la razón social no exista
	if req.RazonSocial != nil && *req.RazonSocial != "" {
		if err := s.validateUniqueRazonSocial(*req.RazonSocial, 0); err != nil {
			return nil, err
		}
	}

	// VALIDACIÓN 4: Verificar consistencia entre tipo de identificación y datos
	if err := s.validateCitizenDataConsistency(req); err != nil {
		return nil, err
	}

	// Crear el modelo desde el DTO
	citizen := s.createCitizenFromRequest(req)

	// Guardar en base de datos
	if err := db.Create(&citizen).Error; err != nil {
		return nil, errors.New("failed to create citizen")
	}

	return s.toCitizenResponse(&citizen), nil
}

// UpdateCitizen actualiza un ciudadano existente
func (s *CitizenService) UpdateCitizen(id uint, req *dto.UpdateCitizenRequest) (*dto.CitizenResponse, error) {
	db := database.GetDB()
	var citizen models.Citizen

	// Obtener ciudadano existente
	if err := db.First(&citizen, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("citizen not found")
		}
		return nil, err
	}

	// Validar cambios únicos si se están modificando
	if req.NumeroIdentificacion != nil && *req.NumeroIdentificacion != citizen.NumeroIdentificacion {
		if err := s.validateUniqueNumeroIdentificacion(*req.NumeroIdentificacion, id); err != nil {
			return nil, err
		}
	}

	if req.Email != nil && *req.Email != citizen.Email {
		if err := s.validateUniqueEmail(*req.Email, id); err != nil {
			return nil, err
		}
	}

	if req.RazonSocial != nil && *req.RazonSocial != "" {
		// Comparar con el valor actual (que puede ser nil)
		currentRazonSocial := ""
		if citizen.RazonSocial != nil {
			currentRazonSocial = *citizen.RazonSocial
		}
		if *req.RazonSocial != currentRazonSocial {
			if err := s.validateUniqueRazonSocial(*req.RazonSocial, id); err != nil {
				return nil, err
			}
		}
	}

	// Aplicar cambios al modelo
	s.applyCitizenUpdates(&citizen, req)

	// Guardar cambios
	if err := db.Save(&citizen).Error; err != nil {
		return nil, errors.New("failed to update citizen")
	}

	return s.toCitizenResponse(&citizen), nil
}

// DeleteCitizen elimina un ciudadano (soft delete)
func (s *CitizenService) DeleteCitizen(id uint) error {
	db := database.GetDB()

	// Verificar que el ciudadano existe
	var citizen models.Citizen
	if err := db.First(&citizen, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("citizen not found")
		}
		return err
	}

	// Soft delete - GORM automáticamente maneja el DeletedAt
	return db.Delete(&citizen).Error
}

// --- MÉTODOS DE VALIDACIÓN PRIVADOS ---

// validateUniqueNumeroIdentificacion verifica que el número de identificación sea único
func (s *CitizenService) validateUniqueNumeroIdentificacion(numero string, excludeID uint) error {
	db := database.GetDB()
	var count int64

	query := db.Model(&models.Citizen{}).Where("numero_identificacion = ?", numero)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errors.New("error validating identification number")
	}

	if count > 0 {
		return errors.New("identification number already exists")
	}

	return nil
}

// validateUniqueEmail verifica que el email sea único
func (s *CitizenService) validateUniqueEmail(email string, excludeID uint) error {
	db := database.GetDB()
	var count int64

	query := db.Model(&models.Citizen{}).Where("email = ?", email)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errors.New("error validating email")
	}

	if count > 0 {
		return errors.New("email already exists")
	}

	return nil
}

// validateUniqueRazonSocial verifica que la razón social sea única
func (s *CitizenService) validateUniqueRazonSocial(razonSocial string, excludeID uint) error {
	db := database.GetDB()
	var count int64

	query := db.Model(&models.Citizen{}).Where("razon_social = ?", razonSocial)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errors.New("error validating razon social")
	}

	if count > 0 {
		return errors.New("razon social already exists")
	}

	return nil
}

// validateCitizenDataConsistency verifica que los datos sean consistentes
// Por ejemplo: si es RUC (04) debe tener razón social, si es cédula (05) debe tener nombre
func (s *CitizenService) validateCitizenDataConsistency(req *dto.CreateCitizenRequest) error {
	switch req.TipoIdentificacion {
	case "04": // RUC - debe ser empresa
		if req.RazonSocial == nil || *req.RazonSocial == "" {
			return errors.New("RUC requires razon_social")
		}
	case "05", "06": // Cédula o Pasaporte - debe ser persona natural
		if req.Nombre == nil || *req.Nombre == "" {
			return errors.New("natural person requires nombre")
		}
	}
	return nil
}

// --- MÉTODOS DE CONVERSIÓN PRIVADOS ---

// createCitizenFromRequest convierte un DTO de creación a modelo
func (s *CitizenService) createCitizenFromRequest(req *dto.CreateCitizenRequest) models.Citizen {
	return models.Citizen{
		// Identificación
		NumeroIdentificacion: req.NumeroIdentificacion,
		TipoIdentificacion:   req.TipoIdentificacion,
		
		// Contacto
		Email:              req.Email,
		Celular:            req.Celular,
		Convencional:       req.Convencional,
		DireccionPrincipal: req.DireccionPrincipal,
		Pais:               req.Pais,
		Provincia:          req.Provincia,
		Ciudad:             req.Ciudad,

		// Persona Natural
		Nombre:          req.Nombre,
		FechaNacimiento: req.FechaNacimiento,
		Nacionalidad:    req.Nacionalidad,
		EstadoCivil:     req.EstadoCivil,
		Genero:          req.Genero,

		// Empresa
		RazonSocial:           req.RazonSocial,
		NombreComercial:       req.NombreComercial,
		TipoEmpresa:           req.TipoEmpresa,
		RepresentantesLegales: req.RepresentantesLegales,

		// Tributario
		TipoContribuyente:           req.TipoContribuyente,
		EstadoContribuyente:         req.EstadoContribuyente,
		Regimen:                     req.Regimen,
		Categoria:                   req.Categoria,
		ObligadoContabilidad:        req.ObligadoContabilidad,
		AgenteRetencion:             req.AgenteRetencion,
		ContribuyenteEspecial:       req.ContribuyenteEspecial,
		ActividadEconomicaPrincipal: req.ActividadEconomicaPrincipal,
		Sucursales:                  req.Sucursales,

		// Metadatos
		MotivoCancelacionSuspension: req.MotivoCancelacionSuspension,
	}
}

// applyCitizenUpdates aplica actualizaciones parciales al modelo
func (s *CitizenService) applyCitizenUpdates(citizen *models.Citizen, req *dto.UpdateCitizenRequest) {
	// Solo actualizar campos que no son nulos en el request
	if req.NumeroIdentificacion != nil {
		citizen.NumeroIdentificacion = *req.NumeroIdentificacion
	}
	if req.TipoIdentificacion != nil {
		citizen.TipoIdentificacion = *req.TipoIdentificacion
	}
	if req.Email != nil {
		citizen.Email = *req.Email
	}
	if req.Celular != nil {
		citizen.Celular = *req.Celular
	}
	if req.Convencional != nil {
		citizen.Convencional = *req.Convencional
	}
	if req.DireccionPrincipal != nil {
		citizen.DireccionPrincipal = *req.DireccionPrincipal
	}
	if req.Pais != nil {
		citizen.Pais = *req.Pais
	}
	if req.Provincia != nil {
		citizen.Provincia = *req.Provincia
	}
	if req.Ciudad != nil {
		citizen.Ciudad = *req.Ciudad
	}

	// Persona Natural
	if req.Nombre != nil {
		citizen.Nombre = req.Nombre
	}
	if req.FechaNacimiento != nil {
		citizen.FechaNacimiento = req.FechaNacimiento
	}
	if req.Nacionalidad != nil {
		citizen.Nacionalidad = req.Nacionalidad
	}
	if req.EstadoCivil != nil {
		citizen.EstadoCivil = req.EstadoCivil
	}
	if req.Genero != nil {
		citizen.Genero = req.Genero
	}

	// Empresa
	if req.RazonSocial != nil {
		citizen.RazonSocial = req.RazonSocial
	}
	if req.NombreComercial != nil {
		citizen.NombreComercial = req.NombreComercial
	}
	if req.TipoEmpresa != nil {
		citizen.TipoEmpresa = req.TipoEmpresa
	}
	if req.RepresentantesLegales != nil {
		citizen.RepresentantesLegales = req.RepresentantesLegales
	}

	// Tributario
	if req.TipoContribuyente != nil {
		citizen.TipoContribuyente = *req.TipoContribuyente
	}
	if req.EstadoContribuyente != nil {
		citizen.EstadoContribuyente = *req.EstadoContribuyente
	}
	if req.Regimen != nil {
		citizen.Regimen = *req.Regimen
	}
	if req.Categoria != nil {
		citizen.Categoria = *req.Categoria
	}
	if req.ObligadoContabilidad != nil {
		citizen.ObligadoContabilidad = *req.ObligadoContabilidad
	}
	if req.AgenteRetencion != nil {
		citizen.AgenteRetencion = req.AgenteRetencion
	}
	if req.ContribuyenteEspecial != nil {
		citizen.ContribuyenteEspecial = req.ContribuyenteEspecial
	}
	if req.ActividadEconomicaPrincipal != nil {
		citizen.ActividadEconomicaPrincipal = *req.ActividadEconomicaPrincipal
	}
	if req.Sucursales != nil {
		citizen.Sucursales = req.Sucursales
	}

	// Metadatos
	if req.MotivoCancelacionSuspension != nil {
		citizen.MotivoCancelacionSuspension = *req.MotivoCancelacionSuspension
	}
}

// toCitizenResponse convierte un modelo a DTO de respuesta
func (s *CitizenService) toCitizenResponse(citizen *models.Citizen) *dto.CitizenResponse {
	response := &dto.CitizenResponse{
		ID: citizen.ID,

		// Identificación
		NumeroIdentificacion: citizen.NumeroIdentificacion,
		TipoIdentificacion:   citizen.TipoIdentificacion,
		
		// Contacto
		Email:              citizen.Email,
		Celular:            citizen.Celular,
		Convencional:       citizen.Convencional,
		DireccionPrincipal: citizen.DireccionPrincipal,
		Pais:               citizen.Pais,
		Provincia:          citizen.Provincia,
		Ciudad:             citizen.Ciudad,

		// Persona Natural
		Nombre:          citizen.Nombre,
		FechaNacimiento: citizen.FechaNacimiento,
		Nacionalidad:    citizen.Nacionalidad,
		EstadoCivil:     citizen.EstadoCivil,
		Genero:          citizen.Genero,

		// Empresa
		RazonSocial:           citizen.RazonSocial,
		NombreComercial:       citizen.NombreComercial,
		TipoEmpresa:           citizen.TipoEmpresa,
		RepresentantesLegales: citizen.RepresentantesLegales,

		// Tributario
		TipoContribuyente:           citizen.TipoContribuyente,
		EstadoContribuyente:         citizen.EstadoContribuyente,
		Regimen:                     citizen.Regimen,
		Categoria:                   citizen.Categoria,
		ObligadoContabilidad:        citizen.ObligadoContabilidad,
		AgenteRetencion:             citizen.AgenteRetencion,
		ContribuyenteEspecial:       citizen.ContribuyenteEspecial,
		ActividadEconomicaPrincipal: citizen.ActividadEconomicaPrincipal,
		Sucursales:                  citizen.Sucursales,

		// Metadatos
		MotivoCancelacionSuspension: citizen.MotivoCancelacionSuspension,
		CreatedAt:                   citizen.CreatedAt,
		UpdatedAt:                   citizen.UpdatedAt,
	}

	// Calcular edad si hay fecha de nacimiento
	if citizen.FechaNacimiento != nil {
		age := calculateAge(*citizen.FechaNacimiento)
		response.Edad = &age
	}

	return response
}

// calculateAge calcula la edad basada en la fecha de nacimiento
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	
	// Ajustar si aún no ha cumplido años este año
	if now.YearDay() < birthDate.YearDay() {
		age--
	}
	
	return age
}