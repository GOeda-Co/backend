package memory

import (
	"errors"
	"fmt"
	"sync"

	"repeatro/src/user/pkg/model"

	"github.com/google/uuid"
)

type Repository struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*model.User
}

func NewInMemoryUserRepo() *Repository {
	return &Repository{
		users: make(map[uuid.UUID]*model.User),
	}
}

func (r *Repository) CreateUser(user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.UserId == uuid.Nil {
		user.UserId = uuid.New()
	}
	// Optional: Check for duplicate email
	for _, u := range r.users {
		if u.Email == user.Email {
			return errors.New("email already exists")
		}
	}

	r.users[user.UserId] = user
	fmt.Println("TOO") // matching original side-effect
	return nil
}

func (r *Repository) ReadUser(userID uuid.UUID) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *Repository) ReadUserByEmail(email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *Repository) ReadAllUsers() ([]model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []model.User
	for _, u := range r.users {
		result = append(result, *u)
	}
	return result, nil
}
