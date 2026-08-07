package user_repository

import (
	"database/sql"
	"regexp"
	"testing"

	"VSRT-Lang/internal/user"

	"github.com/DATA-DOG/go-sqlmock"
)

type sqlmockRows = sqlmock.Rows

func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock db: %v", err)
	}
	return db, mock
}

func TestRepository_Create_SetsIDWhenEmpty(t *testing.T) {
	db, mock := setupDB(t)

	user := user.NewUser("demo", "demo@example.com", "secret")

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(sqlmock.AnyArg(), user.Username, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	defer db.Close()

	repo := New(db)

	if err := repo.Create(user); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID == "" {
		t.Fatal("expected user ID to be set")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_Create_UsesProvidedID(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	user := user.NewUser("demo", "demo@example.com", "secret")

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users (id, username, email, password, created_at)
		VALUES ($1, $2, $3, $4, $5)`)).
		WithArgs(sqlmock.AnyArg(), user.Username, user.Email, user.Password, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := New(db)

	if err := repo.Create(user); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByEmail_ReturnsUser(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	expected := user.NewUser("demo", "demo@example.com", "secret")

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at"}).
		AddRow(expected.ID, expected.Username, expected.Email, expected.Password, expected.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE email = $1`)).
		WithArgs(expected.Email).
		WillReturnRows(rows)

	repo := New(db)
	got, err := repo.FindByEmail(expected.Email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != expected.ID || got.Username != expected.Username || got.Email != expected.Email || got.Password != expected.Password || !got.CreatedAt.Equal(expected.CreatedAt) {
		t.Fatalf("unexpected user returned: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByEmail_NotFound(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE email = $1`)).
		WithArgs("missing@example.com").
		WillReturnError(sql.ErrNoRows)

	repo := New(db)
	_, err := repo.FindByEmail("missing@example.com")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByUsername_ReturnsUser(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	expected := user.NewUser("demo", "demo@example.com", "secret")

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at"}).
		AddRow(expected.ID, expected.Username, expected.Email, expected.Password, expected.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE username = $1`)).
		WithArgs(expected.Username).
		WillReturnRows(rows)

	repo := New(db)
	got, err := repo.FindByUsername(expected.Username)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != expected.ID || got.Username != expected.Username || got.Email != expected.Email || got.Password != expected.Password || !got.CreatedAt.Equal(expected.CreatedAt) {
		t.Fatalf("unexpected user returned: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByUsername_NotFound(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE username = $1`)).
		WithArgs("missing-username").
		WillReturnError(sql.ErrNoRows)

	repo := New(db)
	_, err := repo.FindByUsername("missing-username")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByID_ReturnsUser(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	expected := user.NewUser("demo", "demo@example.com", "secret")

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at"}).
		AddRow(expected.ID, expected.Username, expected.Email, expected.Password, expected.CreatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE id = $1`)).
		WithArgs(expected.ID).
		WillReturnRows(rows)

	repo := New(db)
	got, err := repo.FindByID(expected.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.ID != expected.ID || got.Username != expected.Username || got.Email != expected.Email || got.Password != expected.Password || !got.CreatedAt.Equal(expected.CreatedAt) {
		t.Fatalf("unexpected user returned: %#v", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestRepository_FindByID_NotFound(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password, created_at
         FROM users
         WHERE id = $1`)).
		WithArgs("missing-id").
		WillReturnError(sql.ErrNoRows)

	repo := New(db)
	_, err := repo.FindByID("missing-id")
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
