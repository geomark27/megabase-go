// internal/database/seeders/userSeeder.go
package seeders

import (
    "log"

    "megabaseGo/internal/models"
    "megabaseGo/internal/utils"

    "gorm.io/gorm"
)

type UserSeeder struct {
    Hasher utils.PasswordHasher
}

// NewUserSeeder inyecta el hasher al seeder
func NewUserSeeder(hasher utils.PasswordHasher) *UserSeeder {
    return &UserSeeder{Hasher: hasher}
}

func (s *UserSeeder) Run(db *gorm.DB) error {
    // 1) Genera el hash con tu util
    hashedPassword, err := s.Hasher.HashPassword("admin123")
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return err
    }

    // 2) Verifica existencia
    var existing models.User
    res := db.Where("user_name = ?", "admin").First(&existing)
    if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
        log.Printf("Error revisando existencia de usuario admin: %v", res.Error)
        return res.Error
    }

    // 3) Si no existe, lo crea (GORM setea CreatedAt/UpdatedAt automáticamente)
    if res.Error == gorm.ErrRecordNotFound {
        admin := models.User{
            Name:     "Admin",
            UserName: "admin",
            Email:    "admin@admin.com",
            Password: hashedPassword,
            RoleID:   1,
            IsActive: true,
            // LastLoginAt queda en cero, GORM lo manejará si tienes hooks
        }
        if err := db.Create(&admin).Error; err != nil {
            log.Printf("Error creando el usuario admin: %v", err)
            return err
        }
        log.Println("Creacion de usuario admin exitosa")
    } else {
        log.Println("El usuario admin ya existe, se omite la creacion")
    }
    return nil
}
