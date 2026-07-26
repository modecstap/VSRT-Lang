package refresh_token_repository

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

func (r *Repository) Save(token *auth.RefreshToken) error {
	if token.ID == "" {
		token.ID = fmt.Sprintf("refresh-%d", time.Now().UnixNano())
	}

	_, err := r.db.Exec(
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked)
         VALUES ($1, $2, $3, $4, $5, $6)
         ON CONFLICT (token_hash) DO UPDATE SET
             user_id    = excluded.user_id,
             expires_at = excluded.expires_at,
             created_at = excluded.created_at,
             revoked    = excluded.revoked`,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	return err
}

func (r *Repository) FindByHash(hash string) (*auth.RefreshToken, error) {
	token := &auth.RefreshToken{}
	err := r.db.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, created_at, revoked
         FROM refresh_tokens
         WHERE token_hash = $1`,
		hash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.Revoked,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return token, nil
}

func (r *Repository) RevokeByHash(hash string) error {
	result, err := r.db.Exec(
		`UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`,
		hash,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("refresh token not found")
	}
	return nil
}

func (r *Repository) DeleteByHash(hash string) error {
	result, err := r.db.Exec(
		`DELETE FROM refresh_tokens WHERE token_hash = $1`,
		hash,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("refresh token not found")
	}
	return nil
}
