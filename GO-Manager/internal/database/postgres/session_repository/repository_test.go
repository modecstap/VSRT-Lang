package session_repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock db: %v", err)
	}

	return db, mock
}

func SaveInsertsSessionAndRecords(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	userID := user.UserId("user-1")
	sess := &session.Session{User: userID, Name: "demo"}
	sess.Records = make(map[string]session.Record)
	sess.Records["hello"] = session.Record{
		Phrase:       "hello",
		BaseForm:     "hello",
		Translations: []string{"hola"},
		Synonyms:     []string{"greeting"},
		Antonyms:     []string{"bye"},
		Contexts: []session.Context{{
			Phrase:      "context",
			Translation: "contexto",
		}},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO sessions (user_id, name)
		VALUES ($1, $2)
		RETURNING id
	`)).
		WithArgs(userID, sess.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM records WHERE session_id=$1`)).
		WithArgs(int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(regexp.QuoteMeta(`
        INSERT INTO records(
            session_id,
            phrase,
            base_form,
            translations,
            synonyms,
            antonyms,
            contexts
        )
        VALUES($1,$2,$3,$4,$5,$6,$7)
    `))
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO records(
            session_id,
            phrase,
            base_form,
            translations,
            synonyms,
            antonyms,
            contexts
        )
        VALUES($1,$2,$3,$4,$5,$6,$7)
    `)).
		WithArgs(int64(42), sess.Records["hello"].Phrase, sess.Records["hello"].BaseForm, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := New(db)
	if err := repo.Save(sess); err != nil {
		t.Fatalf("expected save to succeed, got %v", err)
	}

	if sess.ID != 42 {
		t.Fatalf("expected session id to be set to 42, got %d", sess.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func FindByUserReturnsSessionsAndRecords(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	userID := user.UserId("user-1")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			id,
			user_id,
			name
		FROM sessions
		WHERE user_id = $1
	`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(int64(7), userID, "demo"))
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			phrase,
			base_form,
			translations,
			synonyms,
			antonyms,
			contexts
		FROM records
		WHERE session_id = $1
		ORDER BY id
	`)).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"phrase", "base_form", "translations", "synonyms", "antonyms", "contexts"}))

	repo := New(db)
	sessions, err := repo.FindByUser(userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(sessions) != 1 {
		t.Fatalf("expected one session, got %d", len(sessions))
	}

	if sessions[0].User != userID || sessions[0].Name != "demo" {
		t.Fatalf("unexpected session returned: %#v", sessions[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TakeReturnsSession(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	userID := user.UserId("user-1")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			id,
			user_id,
			name
		FROM sessions
		WHERE id = $1
	`)).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(int64(7), userID, "demo"))
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			phrase,
			base_form,
			translations,
			synonyms,
			antonyms,
			contexts
		FROM records
		WHERE session_id = $1
		ORDER BY id
	`)).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"phrase", "base_form", "translations", "synonyms", "antonyms", "contexts"}))

	repo := New(db)
	sess, err := repo.Take(7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sess.ID != 7 || sess.User != userID || sess.Name != "demo" {
		t.Fatalf("unexpected session returned: %#v", sess)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TakeReturnsRecordQueryError(t *testing.T) {
	db, mock := setupDB(t)
	defer db.Close()

	userID := user.UserId("user-1")
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			id,
			user_id,
			name
		FROM sessions
		WHERE id = $1
	`)).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(int64(7), userID, "demo"))
	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT
			phrase,
			base_form,
			translations,
			synonyms,
			antonyms,
			contexts
		FROM records
		WHERE session_id = $1
		ORDER BY id
	`)).
		WithArgs(int64(7)).
		WillReturnError(errors.New("boom"))

	repo := New(db)
	_, err := repo.Take(7)
	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected record query error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
