package card

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	domain "VSRT-Lang/internal/card"
	"VSRT-Lang/internal/session"
)

func TestHandler_UpdateKnowledge(t *testing.T) {
	t.Parallel()

	fails := func(err error) *fakeService {
		return &fakeService{
			updateKnowledge: func(domain.UpdateKnowledgeCommand) error { return err },
		}
	}
	const body = `{"phrase":"Hello","estimation":"EASY"}`

	tests := []struct {
		name          string
		path          string
		body          string
		auth          bool
		service       *fakeService
		wantStatus    int
		wantErrorCode string
	}{
		{
			name: "estimates card",
			path: "/sessions/7/cards",
			body: body,
			auth: true,
			service: &fakeService{
				updateKnowledge: func(cmd domain.UpdateKnowledgeCommand) error {
					want := domain.UpdateKnowledgeCommand{UserId: testUserID, SessionId: 7, Phrase: "Hello", Estimation: domain.EstimationEasy}
					if cmd != want {
						return errors.New("unexpected update knowledge command")
					}
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{name: "rejects unknown route", path: "/sessions/7/cards/x", body: body, auth: true, service: &fakeService{}, wantStatus: http.StatusNotFound, wantErrorCode: "not_found"},
		{name: "rejects invalid session id", path: "/sessions/abc/cards", body: body, auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects invalid json", path: "/sessions/7/cards", body: `{`, auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects blank phrase", path: "/sessions/7/cards", body: `{"phrase":"  ","estimation":"EASY"}`, auth: true, service: &fakeService{}, wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "maps unknown estimation", path: "/sessions/7/cards", body: body, auth: true, service: fails(domain.ErrUnknownEstimation), wantStatus: http.StatusBadRequest, wantErrorCode: "invalid_request"},
		{name: "rejects missing user context", path: "/sessions/7/cards", body: body, service: &fakeService{}, wantStatus: http.StatusUnauthorized, wantErrorCode: "unauthorized"},
		{name: "maps missing session", path: "/sessions/7/cards", body: body, auth: true, service: fails(session.ErrSessionNotFound), wantStatus: http.StatusNotFound, wantErrorCode: "session_not_found"},
		{name: "maps unauthorized as not found", path: "/sessions/7/cards", body: body, auth: true, service: fails(session.ErrUnauthorized), wantStatus: http.StatusNotFound, wantErrorCode: "session_not_found"},
		{name: "maps missing card", path: "/sessions/7/cards", body: body, auth: true, service: fails(domain.ErrCardNotFound), wantStatus: http.StatusNotFound, wantErrorCode: "card_not_found"},
		{name: "maps other service errors", path: "/sessions/7/cards", body: body, auth: true, service: fails(errors.New("db down")), wantStatus: http.StatusInternalServerError, wantErrorCode: "card_estimate_failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.service)
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))

			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, testJWT(t), testUserID, h.UpdateKnowledge, req)
			} else {
				rw = httptest.NewRecorder()
				h.UpdateKnowledge(rw, req)
			}

			if rw.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.wantStatus, rw.Body.String())
			}
			if tt.wantStatus == http.StatusNoContent && rw.Body.Len() != 0 {
				t.Fatalf("body = %q, want empty", rw.Body.String())
			}
			if tt.wantErrorCode != "" {
				got := decodeJSONError(t, rw)
				if got.Error.Code != tt.wantErrorCode {
					t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantErrorCode)
				}
			}
		})
	}
}
