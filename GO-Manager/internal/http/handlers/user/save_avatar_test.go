package user

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/http/handlers"
	domain "VSRT-Lang/internal/user"
)

type recordingUserRepo struct {
	lastID     domain.UserId
	lastAvatar domain.Avatar
	saveErr    error
}

func (r *recordingUserRepo) Create(*domain.User) error                   { return nil }
func (r *recordingUserRepo) FindByEmail(string) (*domain.User, error)    { return nil, nil }
func (r *recordingUserRepo) FindByUsername(string) (*domain.User, error) { return nil, nil }
func (r *recordingUserRepo) FindByID(string) (*domain.User, error)       { return nil, nil }
func (r *recordingUserRepo) GetAvatar(domain.UserId) (domain.Avatar, error) {
	return domain.Avatar{}, domain.ErrNoAvatar
}
func (r *recordingUserRepo) SaveAvatar(id domain.UserId, avatar domain.Avatar) error {
	r.lastID = id
	r.lastAvatar = avatar
	return r.saveErr
}

func pngBytes(t *testing.T, w, h int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, color.RGBA{G: 255, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("png.Encode: %v", err)
	}
	return buf.Bytes()
}

func gifBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	var buf bytes.Buffer
	if err := gif.Encode(&buf, img, nil); err != nil {
		t.Fatalf("gif.Encode: %v", err)
	}
	return buf.Bytes()
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

	repo := &recordingUserRepo{}
	h := NewHandler(repo, nil)
	raw := pngBytes(t, 16, 16)
	rw := serveAuthed(t, h.SaveAvatar, avatarRequest(t, "avatar", "a.png", raw), "user-1")

	if rw.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204; body=%s", rw.Code, rw.Body.String())
	}
	if rw.Header().Get("Content-Type") != "" {
		t.Fatalf("Content-Type = %q, want empty", rw.Header().Get("Content-Type"))
	}
	if repo.lastID != "user-1" {
		t.Fatalf("saved id = %q, want user-1", repo.lastID)
	}
	want, err := domain.PrepareAvatar(raw)
	if err != nil {
		t.Fatalf("PrepareAvatar: %v", err)
	}
	if repo.lastAvatar.MediaType != "image/png" || !bytes.Equal(repo.lastAvatar.Bytes, want.Bytes) {
		t.Fatalf("stored avatar mismatch")
	}
}

func TestHandler_SaveAvatar_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		auth  bool
		field string
		file  []byte
		want  int
		code  string
	}{
		{name: "missing field", auth: true, field: "file", file: pngBytes(t, 8, 8), want: http.StatusBadRequest, code: "avatar_missing"},
		{name: "gif", auth: true, field: "avatar", file: gifBytes(t), want: http.StatusBadRequest, code: "invalid_avatar_type"},
		{name: "oversized", auth: true, field: "avatar", file: make([]byte, domain.MaxAvatarBytes+1), want: http.StatusBadRequest, code: "avatar_too_large"},
		{name: "513px", auth: true, field: "avatar", file: pngBytes(t, 513, 1), want: http.StatusBadRequest, code: "avatar_dimensions_invalid"},
		{name: "no context", auth: false, field: "avatar", file: pngBytes(t, 8, 8), want: http.StatusUnauthorized, code: "unauthorized"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewHandler(&recordingUserRepo{}, nil)
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
		})
	}
}

func TestHandler_SaveAvatar_IsolatedByToken(t *testing.T) {
	t.Parallel()

	repo := memory.NewUserRepository()
	a := domain.NewUser("alice", "alice@example.com", "secret")
	a.ID = "user-a"
	b := domain.NewUser("bob", "bob@example.com", "secret")
	b.ID = "user-b"
	if err := repo.Create(a); err != nil {
		t.Fatalf("create a: %v", err)
	}
	if err := repo.Create(b); err != nil {
		t.Fatalf("create b: %v", err)
	}

	h := NewHandler(repo, nil)
	pngA := pngBytes(t, 8, 8)
	pngB := pngBytes(t, 16, 16)

	rwA := serveAuthed(t, h.SaveAvatar, avatarRequest(t, "avatar", "a.png", pngA), "user-a")
	if rwA.Code != http.StatusNoContent {
		t.Fatalf("A status = %d, body=%s", rwA.Code, rwA.Body.String())
	}
	rwB := serveAuthed(t, h.SaveAvatar, avatarRequest(t, "avatar", "b.png", pngB), "user-b")
	if rwB.Code != http.StatusNoContent {
		t.Fatalf("B status = %d, body=%s", rwB.Code, rwB.Body.String())
	}

	gotA, err := repo.GetAvatar("user-a")
	if err != nil {
		t.Fatalf("GetAvatar A: %v", err)
	}
	gotB, err := repo.GetAvatar("user-b")
	if err != nil {
		t.Fatalf("GetAvatar B: %v", err)
	}
	if bytes.Equal(gotA.Bytes, gotB.Bytes) {
		t.Fatal("B save overwrote A avatar")
	}
	wantA, err := domain.PrepareAvatar(pngA)
	if err != nil {
		t.Fatalf("PrepareAvatar A: %v", err)
	}
	if !bytes.Equal(gotA.Bytes, wantA.Bytes) {
		t.Fatal("A avatar changed after B save")
	}
}
