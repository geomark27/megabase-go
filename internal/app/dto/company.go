package dto

// CreateCompanyRequest estructura para crear una compañía
type CreateCompanyRequest struct {
	Name   		string `json:"name" binding:"required,max=100"`
	Host     	string `json:"host" binding:"required,max=100"`
	Database 	string `json:"database" binding:"required,max=100"`
	User		string `json:"user" binding:"required,max=100"`
	Password 	string `json:"password" binding:"required,max=100"`
	IsActive 	bool  `json:"is_active"`
}

// UpdateCompanyRequest estructura para actualizar una compañía
type UpdateCompanyRequest struct {
	Name   		*string `json:"name,omitempty" binding:"omitempty,max=100"`
	Host     	*string `json:"host,omitempty" binding:"omitempty,max=100"`
	Database 	*string `json:"database,omitempty" binding:"omitempty,max=100"`
	User     	*string `json:"user,omitempty" binding:"omitempty,max=100"`
	Password 	*string `json:"password,omitempty" binding:"omitempty,max=100"`
	IsActive 	*bool   `json:"is_active,omitempty"`
}

// CompanyResponse estructura para responder datos de una compañía
type CompanyResponse struct {
	ID        	uint        `json:"id"`
	Name    	string      `json:"name"`
	Host      	string      `json:"host"`
	Database  	string      `json:"database"`
	User      	string      `json:"user"`
	IsActive	bool        `json:"is_active"`
	CreatedAt 	interface{} `json:"created_at"`
	UpdatedAt 	interface{} `json:"updated_at"`
}

type CompanySearchFilters struct {
	Name 		*string `form:"name" binding:"omitempty,max=100"`
	IsActive 	*bool   `form:"is_active" binding:"omitempty,oneof=true false"`
}
