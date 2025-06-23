package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId           uuid.UUID `gorm:"type:uuid;primaryKey" json:"user_id" bson:"user_id"`
	Email            string    `gorm:"type:varchar(255);not null;default:null" json:"email" bson:"email"`
	HashedPassword   string    `gorm:"column:hashed_password;type:varchar(255);not null" json:"hashed_password" bson:"hashed_password"`
	RegistrationDate time.Time `gorm:"autoCreateTime" json:"registration_date" bson:"registration_date"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserId = uuid.New()
	return nil
}
