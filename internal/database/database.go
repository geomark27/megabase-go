package database

import (
	"log"
	"megabaseGo/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB inicializa la conexión a la base de datos PostgreSQL
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDBConnectionString()
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Conexión a PostgreSQL establecida exitosamente")
	return db, nil
}

// GetDB devuelve la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}

// CloseDB cierra la conexión a la base de datos
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error al obtener la instancia de la base de datos: %v", err)
		return
	}
	sqlDB.Close()
}
