package services

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"repeatro/src/user/internal/repository/postgresql"
	"repeatro/src/user/pkg/model"
	"repeatro/src/user/pkg/scheme"
	"repeatro/internal/security"

)

type userRepository interface {
	CreateUser(user *model.User)  error
	ReadUser(user_id uuid.UUID) (*model.User, error)
	ReadAllUsers() ([]model.User, error)
	ReadUserByEmail(email string) (*model.User, error)
}

type Service struct {
	userRepository *postgresql.Repository
	security       *security.Security
}

func CreateNewService(userRepository *postgresql.Repository, security *security.Security) *Service {
	return &Service{
		userRepository: userRepository,
		security:       security,
	}
}

// type ServiceInterface interface {
// 	FindUser(user_id uuid.UUID) (*model.User, error)
// 	CreateUser(user *model.User) (uuid.UUID, error)
// 	FindAllUsers() ([]model.User, error)
// 	GetUserIdByEmail(email string) (uuid.UUID, error)
// 	GetUserByEmail(email string) (*model.User, error)
// 	Register(userRegister schemes.AuthUser) (string, error)
// 	Login(userLogin schemes.AuthUser) (string, error)
// }

func (us *Service) FindUser(user_id uuid.UUID) (*model.User, error) {
	return us.userRepository.ReadUser(user_id)
}

func (us *Service) GetUserIdByEmail(email string) (uuid.UUID, error) {
	user, err := us.userRepository.ReadUserByEmail(email)
	if err != nil {
		return uuid.UUID{}, err
	}
	if reflect.DeepEqual(user, &model.User{}) {
		return uuid.UUID{}, fmt.Errorf("not found user")
	}
	return user.UserId, nil
}

func (us *Service) GetUserByEmail(email string) (*model.User, error) {
	user, err := us.userRepository.ReadUserByEmail(email)
	if err != nil {
		return &model.User{}, err
	}
	if reflect.DeepEqual(user, &model.User{}) {
		return &model.User{}, fmt.Errorf("not found user")
	}
	return user, nil
}

func (us *Service) CreateUser(user *model.User) (uuid.UUID, error) {
	userInDB, err := us.FindUser(user.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	if !reflect.DeepEqual(userInDB, &model.User{}) {
		return userInDB.UserId, nil
	}

	err = us.userRepository.CreateUser(user)
	if err != nil {
		return uuid.UUID{}, err
	}

	return user.UserId, nil
}

func (us *Service) FindAllUsers() ([]model.User, error) {
	return us.userRepository.ReadAllUsers()
}

func (us *Service) Register(userRegister schemes.AuthUser) (string, error) {
	userInDB, err := us.userRepository.ReadUserByEmail(userRegister.Email)
	if err != nil {
		if err.Error() != "record not found" {
			return "", err
		}
	}

	if userInDB != nil {
		return "", fmt.Errorf("cannot register user with same email twice")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost) // Use the default cost factor
	if err != nil {
		return "", err
	}

	user := model.User{
		Email:          userRegister.Email,
		HashedPassword: string(hashedPassword),
	}

	user_id, err := us.CreateUser(&user)
	if err != nil {
		return "", err
	}

	token, err := us.security.EncodeString(user.HashedPassword, user_id)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (us *Service) Login(userLogin schemes.AuthUser) (string, error) {
	// want to check that users exists and return user
	user, err := us.GetUserByEmail(userLogin.Email)
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userLogin.Password)) != nil {
		return "", err
	}

	// want to encode string and return token
	token, err := us.security.EncodeString(user.HashedPassword, user.UserId)
	if err != nil {
		return "", err
	}

	return token, nil
}
