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
	ErrNotFound     = errors.New("user not found")
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type UserId string

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
	Avatar    Avatar
}

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
	FindByID(id string) (*User, error)
	Save(user *User) error
	SaveAvatar(id UserId, avatar Avatar) error
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

func (u *User) SetAvatar(raw []byte) error {
	avatar, err := prepareAvatar(raw)
	if err != nil {
		return err
	}
	u.Avatar = avatar
	return nil
}

func (u *User) SetPassword(plain string) error {
	if !ValidatePassword(plain) {
		return ErrWeakPassword
	}
	hashed, err := HashPassword(plain)
	if err != nil {
		return err
	}
	u.Password = hashed
	return nil
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
