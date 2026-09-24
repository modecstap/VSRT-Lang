package mail

import (
	"testing"
)

func TestRenderPasswordReset(t *testing.T) {
	msg, err := RenderPasswordReset("http://localhost/reset-password?token=abc")
	if err != nil {
		t.Fatalf("RenderPasswordReset: %v", err)
	}
	if msg.ContentType != "text/plain" {
		t.Fatalf("ContentType = %q", msg.ContentType)
	}
	if msg.Subject != "Reset your password" {
		t.Fatalf("Subject = %q", msg.Subject)
	}
	wantBody := "Set a new password: http://localhost/reset-password?token=abc. The link works for 15 minutes."
	if msg.Body != wantBody {
		t.Fatalf("Body = %q, want %q", msg.Body, wantBody)
	}
}
