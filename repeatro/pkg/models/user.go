package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserId           uuid.UUID `json:"user_id" bson:"user_id"`
	Email            string    `json:"email" bson:"email"`
	HashedPassword   string    `json:"hashed_password" bson:"hashed_password"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date"`
}
