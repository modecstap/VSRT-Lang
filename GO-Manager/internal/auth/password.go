package auth

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidEmail = errors.New("invalid email")
	errWeakPassword = errors.New("password must be at least 8 characters and include uppercase, lowercase and digit")
)

var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errWeakPassword
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

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}
