// internal/utils/password_hasher.go
package utils

import "golang.org/x/crypto/bcrypt"

// PasswordHasher abstrae el hashing de contraseñas
type PasswordHasher interface {
    HashPassword(password string) (string, error)
    ComparePassword(hashedPassword, password string) error
}

// BcryptHasher implementa PasswordHasher usando bcrypt
type BcryptHasher struct{}

// NewBcryptHasher construye un BcryptHasher
func NewBcryptHasher() *BcryptHasher {
    return &BcryptHasher{}
}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
    bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bs), err
}

func (b *BcryptHasher) ComparePassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
