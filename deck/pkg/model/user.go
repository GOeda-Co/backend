package model

import (
	"time"

	"github.com/google/uuid"
)

// NOTE: I don't use rn gorm.Model due to redundancy of some fields

type User struct {
	UserId           uuid.UUID `json:"user_id" bson:"user_id"`
	Email            string    `json:"email" bson:"email"`
	HashedPassword   string    `json:"hashed_password" bson:"hashed_password"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date"`
}
