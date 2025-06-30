package services

import (
	"fmt"
	"megabaseGo/internal/app/dto"
	app_errors "megabaseGo/internal/app/errors"
	"megabaseGo/internal/database"
	"megabaseGo/internal/models"

	"gorm.io/gorm"
)

type CompanyService struct {
	db *gorm.DB
}

func NewCompanyService() *CompanyService {
	return &CompanyService{
		db: database.GetDB(),
	}
}

func (s *CompanyService) validateUniqueName(name string, excludeID uint) error {
	var existingCompany models.Company
	query := s.db.Model(&models.Company{}).Where("name = ?", name)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	err := query.First(&existingCompany).Error

	if err == nil {
		// Encontramos una compañía, así que devolvemos nuestro error de conflicto tipado.
		return app_errors.NewConflictError(fmt.Sprintf("company name '%s' already exists", name))
	}
	if err != gorm.ErrRecordNotFound {
		// Error inesperado de la base de datos.
		return err
	}
	// gorm.ErrRecordNotFound significa que el nombre está libre, así que no hay error.
	return nil
}

func (s *CompanyService) CreateCompany(req *dto.CreateCompanyRequest) (*dto.CompanyResponse, error) {
	if err := s.validateUniqueName(req.Name, 0); err != nil {
		return nil, err
	}
	company := models.Company{
		Name:     req.Name, Host: req.Host, Database: req.Database,
		User: req.User, Password: req.Password, IsActive: req.IsActive,
	}
	if err := s.db.Create(&company).Error; err != nil {
		return nil, err
	}
	return toCompanyResponse(&company), nil
}

func (s *CompanyService) GetCompanyByID(id uint) (*dto.CompanyResponse, error) {
	var company models.Company
	if err := s.db.First(&company, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Devolvemos nuestro error Not Found tipado.
			return nil, app_errors.NewNotFoundError("company", id)
		}
		return nil, err
	}
	return toCompanyResponse(&company), nil
}

func (s *CompanyService) UpdateCompany(id uint, req *dto.UpdateCompanyRequest) (*dto.CompanyResponse, error) {
	var company models.Company
	if err := s.db.First(&company, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, app_errors.NewNotFoundError("company", id)
		}
		return nil, err
	}
	if req.Name != nil && *req.Name != company.Name {
		if err := s.validateUniqueName(*req.Name, id); err != nil {
			return nil, err
		}
		company.Name = *req.Name
	}
	if req.Host != nil { company.Host = *req.Host }
	if req.Database != nil { company.Database = *req.Database }
	if req.User != nil { company.User = *req.User }
	if req.Password != nil && *req.Password != "" { company.Password = *req.Password }
	if req.IsActive != nil { company.IsActive = *req.IsActive }

	if err := s.db.Save(&company).Error; err != nil {
		return nil, err
	}
	return toCompanyResponse(&company), nil
}

func (s *CompanyService) DeleteCompany(id uint) error {
    var company models.Company
    if err := s.db.First(&company, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return app_errors.NewNotFoundError("company", id)
        }
        return err
    }
	if err := s.db.Delete(&models.Company{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *CompanyService) GetCompanies(filters *dto.CompanySearchFilters) ([]dto.CompanyResponse, error) {
	var companies []models.Company
	query := s.db.Model(&models.Company{})
	if filters.Name != nil && *filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+*filters.Name+"%")
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if err := query.Find(&companies).Error; err != nil {
		return nil, err
	}
	resp := make([]dto.CompanyResponse, 0)
	for _, c := range companies {
		resp = append(resp, *toCompanyResponse(&c))
	}
	return resp, nil
}

func toCompanyResponse(company *models.Company) *dto.CompanyResponse {
	return &dto.CompanyResponse{
		ID:        company.ID, Name: company.Name, Host: company.Host,
		Database:  company.Database, User: company.User, IsActive:  company.IsActive,
		CreatedAt: company.CreatedAt, UpdatedAt: company.UpdatedAt,
	}
}