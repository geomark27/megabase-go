# Mega Base - Backend API

[![Go Version](https://img.shields.io/badge/go-1.23.4-blue.svg)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/gin--gonic-gin-blue)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/gorm-ORM-lightgrey)](https://gorm.io/)

## Descripción / Description

Mega Base es un backend API desarrollado en Go que proporciona una base sólida para aplicaciones web. Incluye autenticación JWT, manejo de base de datos PostgreSQL con GORM, y una estructura organizada siguiendo las mejores prácticas de desarrollo en Go.

Mega Base is a Go-based backend API that provides a solid foundation for web applications. It includes JWT authentication, PostgreSQL database handling with GORM, and an organized structure following Go best practices.

## 🚀 Características / Features

- Autenticación JWT segura / Secure JWT Authentication
- ORM con GORM para PostgreSQL / GORM ORM for PostgreSQL
- Validación de datos / Data validation
- Variables de entorno con godotenv / Environment variables with godotenv
- Migraciones y seeders / Migrations and seeders
- CORS habilitado / CORS enabled
- Generación de UUIDs / UUID generation

## 📦 Dependencias / Dependencies

- Go 1.23.4 o superior / or higher
- PostgreSQL

### Instalación de dependencias / Install dependencies

```bash
go get github.com/gin-gonic/gin         # Framework web Gin
go get gorm.io/gorm                   # ORM GORM
go get gorm.io/driver/postgres        # Driver PostgreSQL para GORM
go get github.com/golang-jwt/jwt/v5    # JWT tokens
go get golang.org/x/crypto/bcrypt     # Hashing de contraseñas
go get github.com/joho/godotenv       # Variables de entorno (.env)

# Dependencias adicionales recomendadas / Recommended additional dependencies
go get github.com/gin-contrib/cors            # CORS middleware para Gin
go get github.com/go-playground/validator/v10  # Validación de datos
go get github.com/google/uuid                 # Generación de UUIDs
```

## ⚙️ Configuración / Configuration

1. Copiar el archivo `.env.example` a `.env` y configurar las variables de entorno necesarias:

```bash
cp .env.example .env
```

2. Configurar las variables de conexión a la base de datos en el archivo `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_NAME=nombre_base_datos
JWT_SECRET=tu_clave_secreta_jwt
```

## 🚀 Iniciar el proyecto / Run the project

### Ejecutar migraciones con seeders / Run migrations with seeders:

```bash
go run cmd/console/main.go migrate --seed
```

### Iniciar el servidor de desarrollo / Start development server:

```bash
go run main.go
```

El servidor estará disponible en `http://localhost:8080`

## 📚 Estructura del proyecto / Project Structure

```
.
├── cmd/                 # Punto de entrada de la aplicación
│   └── console/         # Comandos de consola (migraciones, seeders, etc.)
├── internal/
│   ├── config/         # Configuración de la aplicación
│   ├── controllers/     # Controladores
│   ├── middleware/      # Middlewares
│   ├── models/          # Modelos de datos
│   ├── repositories/    # Lógica de acceso a datos
│   └── routes/          # Definición de rutas
├── migrations/          # Archivos de migración
├── pkg/                 # Paquetes reutilizables
└── .env.example         # Plantilla de variables de entorno
```

## 📄 Licencia / License

Este proyecto está bajo la Licencia MIT. Ver el archivo [LICENSE](LICENSE) para más información.

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
