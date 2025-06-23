package seeders

import (
	"log"
	"megabaseGo/internal/models"

	"gorm.io/gorm"
)

type RoleSeeder struct{}

func (s *RoleSeeder) Run(db *gorm.DB) error {
	adminRole := models.Role{
		Name:        "admin",
		DisplayName: "Administrator",
		Description: "Administrador con acceso completo al sistema",
		IsActive:    true,
	}

	// Create admin role with ID 1
	if err := db.FirstOrCreate(&adminRole, models.Role{Name: "admin"}).Error; err != nil {
		log.Printf("Error creando rol admin: %v", err)
		return err
	}

	log.Println("Creacion de rol admin exitosa")
	return nil
}
