package user_repository

import (
	"database/sql"
)

type dbConn interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

type Repository struct {
	db dbConn
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func NewTx(tx *sql.Tx) *Repository {
	return &Repository{db: tx}
}
