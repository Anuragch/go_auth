package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID           string     `bson:"_id" json:"id"`
	FullName     string     `bson:"name" json:"name"`
	Email        string     `bson:"email" json:"email"`
	PasswordHash string     `bson:"passwordHash" json:"passwordHash"`
	CreatedAt    *time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

type Credentials struct {
	ID           string `bson:"_id" json:"id"`
	PasswordHash string `bson:"passwordHash" json:"passwordHash"`
}

type UserBody struct {
	Email    string `bson:"email" json:"email"`
	FullName string `bson:"name" json:"name"`
}

type Claims struct {
	ID string `bson:"_id" json:"id"`
	jwt.RegisteredClaims
}
type UserRepository interface {
	GetAllUsers() []*User
	GetUser(userName string) (*User, error)
	AddUser(user *User) bool
	GetUserCredential(userName string) (*Credentials, error)
}
