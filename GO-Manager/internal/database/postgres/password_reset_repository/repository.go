package password_reset_repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"VSRT-Lang/internal/database/postgres/refresh_token_repository"
	"VSRT-Lang/internal/database/postgres/user_repository"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/user"
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

func (r *Repository) Save(link *passwordreset.ResetLink) (user.UserId, error) {
	var tokenHash any
	if link.TokenHash != "" {
		tokenHash = link.TokenHash
	}
	var expiresAt any
	if link.ExpiresAt != nil {
		expiresAt = *link.ExpiresAt
	}

	_, err := r.db.Exec(
		`INSERT INTO password_resets (user_id, token_hash, sent_at, expires_at)
         VALUES ($1, $2, $3, $4)
         ON CONFLICT (user_id) DO UPDATE SET
             token_hash = excluded.token_hash,
             sent_at    = excluded.sent_at,
             expires_at = excluded.expires_at`,
		string(link.User),
		tokenHash,
		link.SentAt,
		expiresAt,
	)
	if err != nil {
		return "", err
	}

	slog.Info("password reset saved",
		"user_id", string(link.User),
		"sent_at", link.SentAt,
		"expires_at", link.ExpiresAt,
		"has_hash", link.TokenHash != "",
	)
	return link.User, nil
}

func (r *Repository) Take(id user.UserId) (*passwordreset.ResetLink, error) {
	return r.take(id, false)
}

func (r *Repository) TakeForUpdate(id user.UserId) (*passwordreset.ResetLink, error) {
	return r.take(id, true)
}

func (r *Repository) take(id user.UserId, forUpdate bool) (*passwordreset.ResetLink, error) {
	query := `SELECT user_id, token_hash, sent_at, expires_at FROM password_resets WHERE user_id = $1`
	if forUpdate {
		query += ` FOR UPDATE`
	}

	var (
		userID    string
		tokenHash sql.NullString
		sentAt    time.Time
		expiresAt sql.NullTime
	)
	err := r.db.QueryRow(query, string(id)).Scan(&userID, &tokenHash, &sentAt, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	link := &passwordreset.ResetLink{
		User:   user.UserId(userID),
		SentAt: sentAt,
	}
	if tokenHash.Valid {
		link.TokenHash = tokenHash.String
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		link.ExpiresAt = &t
	}

	slog.Info("password reset taken",
		"user_id", userID,
		"sent_at", sentAt,
		"expires_at", link.ExpiresAt,
		"has_hash", link.TokenHash != "",
	)
	return link, nil
}

func (r *Repository) TakeByHash(hash string) (*passwordreset.ResetLink, error) {
	var (
		userID    string
		tokenHash sql.NullString
		sentAt    time.Time
		expiresAt sql.NullTime
	)
	err := r.db.QueryRow(
		`SELECT user_id, token_hash, sent_at, expires_at FROM password_resets WHERE token_hash = $1`,
		hash,
	).Scan(&userID, &tokenHash, &sentAt, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	link := &passwordreset.ResetLink{
		User:   user.UserId(userID),
		SentAt: sentAt,
	}
	if tokenHash.Valid {
		link.TokenHash = tokenHash.String
	}
	if expiresAt.Valid {
		t := expiresAt.Time
		link.ExpiresAt = &t
	}

	slog.Info("password reset taken by hash",
		"user_id", userID,
		"sent_at", sentAt,
		"expires_at", link.ExpiresAt,
		"has_hash", link.TokenHash != "",
	)
	return link, nil
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) Within(fn func(passwordreset.TxRepos) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	repos := passwordreset.TxRepos{
		Users:   user_repository.NewTx(tx),
		Links:   NewTx(tx),
		Refresh: refresh_token_repository.NewTx(tx),
	}
	if err := fn(repos); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
