package user_repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"VSRT-Lang/internal/auth"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user *auth.User) error {
	if user.ID == "" {
		user.ID = fmt.Sprintf("user-%d", time.Now().UnixNano())
	}

	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password, created_at)
         VALUES ($1, $2, $3, $4, $5)`,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
	)
	return err
}

func (r *Repository) FindByEmail(email string) (*auth.User, error) {
	user := &auth.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at
         FROM users
         WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByUsername(username string) (*auth.User, error) {
	user := &auth.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at
         FROM users
         WHERE username = $1`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindByID(id string) (*auth.User, error) {
	user := &auth.User{}
	err := r.db.QueryRow(
		`SELECT id, username, email, password, created_at
         FROM users
         WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
