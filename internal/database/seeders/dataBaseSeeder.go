package seeders

import (
	"log"
	"gorm.io/gorm"
)

type DatabaseSeeder struct{}

func (s *DatabaseSeeder) Run(db *gorm.DB) error {
    for _, seeder := range AllSeeders {
        if err := seeder.Run(db); err != nil {
            return err
        }
    }
    log.Println("Ejecutados todos los seeders de la base de datos con éxito")
    return nil
}
