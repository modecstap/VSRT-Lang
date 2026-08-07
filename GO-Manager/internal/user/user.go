package user

import (
	"errors"
	"regexp"
	"time"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail = errors.New("invalid email")
	ErrWeakPassword = errors.New("password must be at least 8 characters and include uppercase, lowercase and digit")
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type UserId string

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByID(id string) (*User, error)
}

func NewUser(username, email, password string) *User {
	return &User{
		ID:        uuid.New().String(),
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", ErrWeakPassword
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}


