package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
	SSLMode    string
}

func LoadConfig() *Config {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró el archivo .env, usando variables de entorno del sistema")
	}

	// Obtener el puerto de la base de datos
	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Error al convertir el puerto de la base de datos: %v", err)
	}

	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "megabase_go"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
	}
}

// GetDBConnectionString devuelve la cadena de conexión para PostgreSQL
func (c *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
