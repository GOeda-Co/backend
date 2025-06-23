package postgresql

import (
	"fmt"
	"log"
	
	"repeatro/src/user/pkg/model"
	

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repository struct {
	db *gorm.DB
}

func NewPostgresRepo(config *viper.Viper, newLogger logger.Interface) *Repository {
	db, err := gorm.Open(postgres.Open(config.GetString("database.connection_string")), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Fatalf("Error during opening database")
	}

	db.AutoMigrate(&model.User{})

	return &Repository{db: db}
}

func (ur Repository) CreateUser(user *model.User) error {
	fmt.Println("TOO")
	return ur.db.Create(user).Error	
}

func (ur Repository) ReadUser(user_id uuid.UUID) (*model.User, error) {
	var user model.User
	err := ur.db.Where("user_id = ?", user_id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur Repository) ReadUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur Repository) ReadAllUsers() ([]model.User, error) {
	var users []model.User
	err := ur.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

