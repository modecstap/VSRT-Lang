package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/http/handlers"
	domain "VSRT-Lang/internal/user"
)

type fakeAvatarService struct {
	id     domain.UserId
	raw    []byte
	err    error
	called int
}

func (f *fakeAvatarService) SaveAvatar(id domain.UserId, raw []byte) error {
	f.called++
	f.id = id
	f.raw = bytes.Clone(raw)
	return f.err
}

func avatarRequest(t *testing.T, field, filename string, data []byte) *http.Request {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	if field != "" {
		part, err := mw.CreateFormFile(field, filename)
		if err != nil {
			t.Fatalf("CreateFormFile: %v", err)
		}
		if _, err := part.Write(data); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}
	if err := mw.Close(); err != nil {
		t.Fatalf("close multipart: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/users/avatar", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func TestHandler_SaveAvatar_Success(t *testing.T) {
	t.Parallel()

	svc := &fakeAvatarService{}
	h := NewHandler(svc, nil)
	raw := []byte("raw-upload")
	rw := serveAuthed(t, h.SaveAvatar, avatarRequest(t, "avatar", "a.png", raw), "user-1")

	if rw.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rw.Code, rw.Body.String())
	}
	if rw.Header().Get("Content-Type") != "" {
		t.Fatalf("Content-Type = %q, want empty", rw.Header().Get("Content-Type"))
	}
	if svc.id != "user-1" {
		t.Fatalf("saved id = %q, want user-1", svc.id)
	}
	if !bytes.Equal(svc.raw, raw) {
		t.Fatalf("saved bytes = %q, want raw upload", svc.raw)
	}
}

func TestHandler_SaveAvatar_Errors(t *testing.T) {
	t.Parallel()

	oversize := make([]byte, domain.MaxAvatarBytes+multipartEnvelopeSlack+1)

	tests := []struct {
		name   string
		auth   bool
		field  string
		file   []byte
		svcErr error
		want   int
		code   string
		called bool
	}{
		{name: "missing field", auth: true, field: "file", file: []byte("x"), want: http.StatusBadRequest, code: "avatar_missing"},
		{name: "no context", auth: false, field: "avatar", file: []byte("x"), want: http.StatusUnauthorized, code: "unauthorized"},
		{name: "body over max reader", auth: true, field: "avatar", file: oversize, want: http.StatusBadRequest, code: "avatar_too_large"},
		{name: "invalid type", auth: true, field: "avatar", file: []byte("x"), svcErr: domain.ErrInvalidAvatarType, want: http.StatusBadRequest, code: "invalid_avatar_type", called: true},
		{name: "dimensions", auth: true, field: "avatar", file: []byte("x"), svcErr: domain.ErrAvatarDimensions, want: http.StatusBadRequest, code: "avatar_dimensions_invalid", called: true},
		{name: "too large", auth: true, field: "avatar", file: []byte("x"), svcErr: domain.ErrAvatarTooLarge, want: http.StatusBadRequest, code: "avatar_too_large", called: true},
		{name: "missing user", auth: true, field: "avatar", file: []byte("x"), svcErr: errors.New("user not found"), want: http.StatusNotFound, code: "user_not_found", called: true},
		{name: "persist failed", auth: true, field: "avatar", file: []byte("x"), svcErr: errors.New("db down"), want: http.StatusInternalServerError, code: "avatar_save_failed", called: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := &fakeAvatarService{err: tt.svcErr}
			h := NewHandler(svc, nil)
			req := avatarRequest(t, tt.field, "a.bin", tt.file)
			var rw *httptest.ResponseRecorder
			if tt.auth {
				rw = serveAuthed(t, h.SaveAvatar, req, "user-1")
			} else {
				rw = httptest.NewRecorder()
				h.SaveAvatar(rw, req)
			}
			if rw.Code != tt.want {
				t.Fatalf("status = %d, want %d; body=%s", rw.Code, tt.want, rw.Body.String())
			}
			var got handlers.ErrorResponse
			if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
				t.Fatalf("decode error: %v", err)
			}
			if got.Error.Code != tt.code {
				t.Fatalf("error code = %q, want %q", got.Error.Code, tt.code)
			}
			if tt.called && svc.called != 1 {
				t.Fatalf("service calls = %d, want 1", svc.called)
			}
			if !tt.called && svc.called != 0 {
				t.Fatalf("service calls = %d, want 0", svc.called)
			}
		})
	}
}
