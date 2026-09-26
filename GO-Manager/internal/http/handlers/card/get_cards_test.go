package card

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

func TestHandler_GetCards(t *testing.T) {
	t.Parallel()

	returns := func(wantLimit int, records []session.Record) *fakeService {
		return &fakeService{
			getCards: func(userID user.UserId, sessionID int64, limit int) ([]session.Record, error) {
				if userID != testUserID || sessionID != 9 || limit != wantLimit {
					return nil, errors.New("unexpected get cards arguments")
				}
				return records, nil
			},
		}
	}
	fails := func(err error) *fakeService {
		return &fakeService{
			getCards: func(user.UserId, int64, int) ([]session.Record, error) { return nil, err },
		}
	}

	tests := []struct {
		name          string
		path          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
		wantBody      string
	}{
		{
			name: "maps records to cards",
			path: "/sessions/9/cards",
			auth: true,
			service: returns(0, []session.Record{
				{
					Phrase:       "hello",
					Translations: []string{"привет"},
					Synonyms:     []string{"hi"},
					Antonyms:     []string{"bye"},
					BaseForm:     "hello",
					Contexts:     []session.Context{{Phrase: "c1", Translation: "t1"}, {Phrase: "c2", Translation: "t2"}},
					Count:        3,
				},
				{Phrase: "bare"},
			}),
			wantStatus: http.StatusOK,
			wantBody: `{"cards":[` +
				`{"front":{"phrase":"hello","base_form":"hello","contexts":["c1","c2"],"synonyms":["hi"]},` +
				`"back":{"translations":["привет"],"context_translations":["t1","t2"]}},` +
				`{"front":{"phrase":"bare","base_form":"","contexts":[],"synonyms":[]},` +
				`"back":{"translations":[],"context_translations":[]}}]}`,
		},
		{
			name:       "passes limit and returns empty list",
			path:       "/sessions/9/cards?limit=2",
			auth:       true,
			service:    returns(2, nil),
			wantStatus: http.StatusOK,
			wantBody:   `{"cards":[]}`,
		},
		{name: "rejects unknown route", path: "/sessions/9/cards/x", auth: true, service: &fakeService{}, wantStatus: http.StatusNotFound, wantErrorCode: "not_found"},
		{name: "rejects invalid session id", path: "/sessions/abc/cards", auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects zero limit", path: "/sessions/9/cards?limit=0", auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects negative limit", path: "/sessions/9/cards?limit=-1", auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects non-numeric limit", path: "/sessions/9/cards?limit=x", auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects empty limit", path: "/sessions/9/cards?limit=", auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects missing user context", path: "/sessions/9/cards", service: &fakeService{}, wantStatus: http.StatusUnauthorized, wantErrorCode: "unauthorized"},
		{name: "maps missing session", path: "/sessions/9/cards", auth: true, service: fails(session.ErrSessionNotFound), wantStatus: http.StatusNotFound, wantErrorCode: "session_not_found"},
		{name: "maps unauthorized as not found", path: "/sessions/9/cards", auth: true, service: fails(session.ErrUnauthorized), wantStatus: http.StatusNotFound, wantErrorCode: "session_not_found"},
		{name: "maps other service errors", path: "/sessions/9/cards", auth: true, service: fails(errors.New("db down")), wantStatus: http.StatusInternalServerError, wantErrorCode: "cards_get_failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.GetCards, req)
			} else {
				rw = httptest.NewRecorder()
				h.GetCards(rw, req)
			}

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}
			if tt.wantErrorCode != "" {
				got := decodeJSONError(t, rw)
				if got.Error.Code != tt.wantErrorCode {
					t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantErrorCode)
				}
				return
			}
			if got := strings.TrimSpace(rw.Body.String()); got != tt.wantBody {
				t.Fatalf("body = %s, want %s", got, tt.wantBody)
			}
		})
	}
}
