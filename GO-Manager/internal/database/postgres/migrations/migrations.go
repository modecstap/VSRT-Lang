package migrations

import (
	"database/sql"
	"fmt"
	"time"
)

type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_users_sessions_records_refresh_tokens",
		Up: `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    id BIGINT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS records (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    phrase TEXT NOT NULL,
    base_form TEXT,
    translations TEXT[],
    synonyms TEXT[],
    antonyms TEXT[],
    contexts JSONB
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked BOOLEAN NOT NULL DEFAULT FALSE
);
`,
		Down: `
DROP TABLE IF EXISTS refresh_tokens CASCADE;
DROP TABLE IF EXISTS records CASCADE;
DROP TABLE IF EXISTS sessions CASCADE;
DROP TABLE IF EXISTS users CASCADE;
`,
	},
}

const schemaMigrationsTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

func Run(db *sql.DB) error {
	if err := ensureSchemaMigrationsTable(db); err != nil {
		return err
	}

	existing, err := appliedVersions(db)
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if _, applied := existing[migration.Version]; applied {
			continue
		}

		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	return nil
}

func ensureSchemaMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(schemaMigrationsTable)
	return err
}

func appliedVersions(db *sql.DB) (map[int]struct{}, error) {
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]struct{})
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		result[version] = struct{}{}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func applyMigration(db *sql.DB, migration Migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(migration.Up); err != nil {
		tx.Rollback()
		return fmt.Errorf("apply migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	if _, err = tx.Exec(
		`INSERT INTO schema_migrations (version, name, applied_at) VALUES ($1, $2, $3)`,
		migration.Version,
		migration.Name,
		time.Now(),
	); err != nil {
		tx.Rollback()
		return fmt.Errorf("record migration %d (%s): %w", migration.Version, migration.Name, err)
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
