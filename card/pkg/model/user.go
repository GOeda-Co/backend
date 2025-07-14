package model

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	PassHash []byte    `json:"pass_hash"`
	IsAdmin  bool      `json:"is_admin"`
}