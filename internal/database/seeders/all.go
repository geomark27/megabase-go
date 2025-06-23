package seeders

import (
	"megabaseGo/internal/utils"

	"gorm.io/gorm"
)

// Seeder define la interfaz para seeders
type Seeder interface {
    Run(db *gorm.DB) error
}

// AllSeeders contiene todos los seeders para ejecución dinámica
var AllSeeders = []Seeder{
    &RoleSeeder{},
    NewUserSeeder(utils.NewBcryptHasher()),
    // Añade aquí tus nuevos seeders, e.g.: &ProductSeeder{},
}