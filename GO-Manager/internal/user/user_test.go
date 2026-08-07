package user

import (
    "testing"
    "time"

    "github.com/google/uuid"
)

func TestValidateEmail(t *testing.T) {
    cases := []struct {
        name  string
        email string
        want  bool
    }{
        {name: "valid email", email: "user@example.com", want: true},
        {name: "missing at", email: "userexample.com", want: false},
        {name: "missing domain", email: "user@", want: false},
        {name: "space in email", email: "user @example.com", want: false},
        {name: "subdomain", email: "user@mail.example.co.uk", want: true},
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := ValidateEmail(c.email)
            if got != c.want {
                t.Fatalf("ValidateEmail(%q) = %v, want %v", c.email, got, c.want)
            }
        })
    }
}

func TestValidatePassword(t *testing.T) {
    cases := []struct {
        name     string
        password string
        want     bool
    }{
        {name: "valid password", password: "Password1", want: true},
        {name: "no uppercase", password: "password1", want: false},
        {name: "no lowercase", password: "PASSWORD1", want: false},
        {name: "no digit", password: "Password", want: false},
        {name: "too short", password: "P1a", want: false},
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            got := ValidatePassword(c.password)
            if got != c.want {
                t.Fatalf("ValidatePassword(%q) = %v, want %v", c.password, got, c.want)
            }
        })
    }
}

func TestHashAndComparePassword(t *testing.T) {
    password := "Password1"

    hash, err := HashPassword(password)
    if err != nil {
        t.Fatalf("HashPassword returned error: %v", err)
    }

    if hash == password {
        t.Fatal("hashed password should not equal plain password")
    }

    if !ComparePassword(hash, password) {
        t.Fatal("ComparePassword returned false for correct password")
    }

    if ComparePassword(hash, "wrong-password") {
        t.Fatal("ComparePassword returned true for incorrect password")
    }
}

func TestHashPassword_EmptyPassword(t *testing.T) {
    _, err := HashPassword("")
    if err != ErrWeakPassword {
        t.Fatalf("HashPassword(\"\") error = %v, want %v", err, ErrWeakPassword)
    }
}

func TestNewUser(t *testing.T) {
    username := "testuser"
    email := "test@example.com"
    password := "Password1"

    u := NewUser(username, email, password)
    if u == nil {
        t.Fatal("NewUser returned nil")
    }

    if u.Username != username {
        t.Fatalf("Username = %q, want %q", u.Username, username)
    }

    if u.Email != email {
        t.Fatalf("Email = %q, want %q", u.Email, email)
    }

    if u.Password != password {
        t.Fatalf("Password = %q, want %q", u.Password, password)
    }

    if u.CreatedAt.IsZero() {
        t.Fatal("CreatedAt should be set")
    }

    if time.Since(u.CreatedAt) < 0 {
        t.Fatal("CreatedAt should not be in the future")
    }

    if _, err := uuid.Parse(u.ID); err != nil {
        t.Fatalf("ID = %q is not a valid uuid: %v", u.ID, err)
    }
}
