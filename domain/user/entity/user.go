package entity

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// User represents a User schema
type User struct {
	UUID      uuid.UUID `json:"-"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt string    `json:"-"`
	UpdatedAt string    `json:"-"`
	DeletedAt string    `json:"-"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	jwt.StandardClaims
	ID uint
}
