package passwordreset_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"VSRT-Lang/internal/auth"
	"VSRT-Lang/internal/database/memory"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/user"
)

type fakeMailer struct {
	to    string
	link  string
	calls int
	err   error
}

func (m *fakeMailer) Send(to, link string) error {
	m.calls++
	m.to = to
	m.link = link
	return m.err
}

type clock struct {
	now time.Time
}

func (c *clock) Now() time.Time { return c.now }

func newService(t *testing.T, mail *fakeMailer, clk *clock, siteURL string) (
	*passwordreset.Service,
	*memory.UserRepository,
	*memory.PasswordResetRepository,
	*memory.RefreshTokenRepository,
) {
	t.Helper()
	users := memory.NewUserRepository()
	links := memory.NewPasswordResetRepository()
	refresh := memory.NewRefreshTokenRepository()
	tx := memory.NewPasswordResetTransactor(users, links, refresh)
	svc := passwordreset.NewService(users, links, mail, tx, siteURL, clk.Now)
	return svc, users, links, refresh
}

func createUser(t *testing.T, users *memory.UserRepository, email, password string) *user.User {
	t.Helper()
	hashed, err := user.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	u := user.NewUser("demo", email, hashed)
	if err := users.Create(u); err != nil {
		t.Fatalf("Create: %v", err)
	}
	return u
}

func tokenFromLink(link string) string {
	return strings.TrimPrefix(link, "http://localhost/reset-password?token=")
}

func TestRequestReset_UnknownEmail(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, _, _, _ := newService(t, mail, clk, "http://localhost")

	if err := svc.RequestReset("missing@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	if mail.calls != 0 {
		t.Fatal("Send must not be called")
	}
	if err := svc.RequestReset("missing@example.com"); err != nil {
		t.Fatalf("repeat: %v", err)
	}
}

func TestRequestReset_SendsAndStoresHash(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, links, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	if mail.calls != 1 || mail.to != "demo@example.com" {
		t.Fatalf("mail = %+v", mail)
	}
	if !strings.HasPrefix(mail.link, "http://localhost/reset-password?token=") {
		t.Fatalf("link = %q", mail.link)
	}
	raw := tokenFromLink(mail.link)
	stored, err := links.Take(user.UserId(u.ID))
	if err != nil || stored == nil {
		t.Fatalf("Take: %v %#v", err, stored)
	}
	if stored.TokenHash == raw || stored.TokenHash == "" {
		t.Fatal("db must store hash, not raw token")
	}
	if !stored.Active(clk.now) {
		t.Fatal("link must be active")
	}
}

func TestRequestReset_RateLimited(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, links, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("first: %v", err)
	}
	first, _ := links.Take(user.UserId(u.ID))

	if err := svc.RequestReset("demo@example.com"); !errors.Is(err, passwordreset.ErrRateLimited) {
		t.Fatalf("second: %v", err)
	}
	if mail.calls != 1 {
		t.Fatalf("Send calls = %d", mail.calls)
	}
	byHash, _ := links.TakeByHash(first.TokenHash)
	if byHash == nil {
		t.Fatal("previous hash must remain")
	}
}

func TestRequestReset_SendFailedLeavesNoLimit(t *testing.T) {
	mail := &fakeMailer{err: errors.New("smtp down")}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, links, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); !errors.Is(err, passwordreset.ErrSendFailed) {
		t.Fatalf("err = %v", err)
	}
	taken, _ := links.Take(user.UserId(u.ID))
	if taken != nil {
		t.Fatal("no row after failed send")
	}
}

func TestRequestReset_EmptySiteURL(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, _, _ := newService(t, mail, clk, "")
	createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); !errors.Is(err, passwordreset.ErrSendFailed) {
		t.Fatalf("err = %v", err)
	}
	if mail.calls != 0 {
		t.Fatal("Send must not run without valid URL")
	}
}

func TestResetPassword_SuccessAndRevoke(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, links, refresh := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	jwtSvc := auth.NewJWTService("secret")
	authSvc := auth.NewService(users, refresh, jwtSvc)
	pair, err := authSvc.Login("demo", "Password1")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	raw := tokenFromLink(mail.link)
	storedBefore, _ := links.Take(user.UserId(u.ID))

	if err := svc.ResetPassword(raw, "Password1"); err != nil {
		t.Fatalf("same password: %v", err)
	}
	found, _ := users.FindByID(u.ID)
	if !user.ComparePassword(found.Password, "Password1") {
		t.Fatal("same password must remain valid")
	}
	byOld, _ := links.TakeByHash(storedBefore.TokenHash)
	if byOld != nil {
		t.Fatal("consumed hash must not resolve")
	}
	stored, _ := links.Take(user.UserId(u.ID))
	if stored == nil || !stored.CoolingDown(clk.now) {
		t.Fatal("cooldown must remain after consume")
	}
	if err := svc.RequestReset("demo@example.com"); !errors.Is(err, passwordreset.ErrRateLimited) {
		t.Fatalf("rate limit after reset: %v", err)
	}
	if _, err := authSvc.Refresh(pair.AccessToken, pair.RefreshToken); err == nil {
		t.Fatal("refresh must fail after reset")
	}
}

func TestResetPassword_ChangesPassword(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, _, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	if err := svc.ResetPassword(tokenFromLink(mail.link), "NewPass99"); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	found, _ := users.FindByID(u.ID)
	if !user.ComparePassword(found.Password, "NewPass99") {
		t.Fatal("new password must match")
	}
	if user.ComparePassword(found.Password, "Password1") {
		t.Fatal("old password must fail")
	}
}

func TestResetPassword_InvalidCases(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, _, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	raw := tokenFromLink(mail.link)
	before, _ := users.FindByID(u.ID)

	if err := svc.ResetPassword("", "NewPass99"); !errors.Is(err, passwordreset.ErrInvalidLink) {
		t.Fatalf("empty token: %v", err)
	}
	if err := svc.ResetPassword("foreign", "NewPass99"); !errors.Is(err, passwordreset.ErrInvalidLink) {
		t.Fatalf("foreign: %v", err)
	}

	clk.now = clk.now.Add(passwordreset.LinkTTL)
	if err := svc.ResetPassword(raw, "NewPass99"); !errors.Is(err, passwordreset.ErrInvalidLink) {
		t.Fatalf("expired: %v", err)
	}
	after, _ := users.FindByID(u.ID)
	if after.Password != before.Password {
		t.Fatal("password must stay on invalid link")
	}

	clk.now = time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC).Add(passwordreset.SendCooldown)
	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("new request: %v", err)
	}
	rawFresh := tokenFromLink(mail.link)
	if err := svc.ResetPassword(rawFresh, "NewPass99"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	if err := svc.ResetPassword(rawFresh, "OtherPass1"); !errors.Is(err, passwordreset.ErrInvalidLink) {
		t.Fatalf("reuse: %v", err)
	}
}

func TestResetPassword_WeakPassword(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, _, _ := newService(t, mail, clk, "http://localhost")
	u := createUser(t, users, "demo@example.com", "Password1")
	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("RequestReset: %v", err)
	}
	raw := tokenFromLink(mail.link)
	before, _ := users.FindByID(u.ID)

	if err := svc.ResetPassword(raw, "weak"); !errors.Is(err, user.ErrWeakPassword) {
		t.Fatalf("err = %v", err)
	}
	after, _ := users.FindByID(u.ID)
	if after.Password != before.Password {
		t.Fatal("password unchanged on weak")
	}
}

func TestRequestReset_NewMailInvalidatesOld(t *testing.T) {
	mail := &fakeMailer{}
	clk := &clock{now: time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)}
	svc, users, _, _ := newService(t, mail, clk, "http://localhost")
	createUser(t, users, "demo@example.com", "Password1")

	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("first: %v", err)
	}
	oldRaw := tokenFromLink(mail.link)

	clk.now = clk.now.Add(passwordreset.SendCooldown)
	if err := svc.RequestReset("demo@example.com"); err != nil {
		t.Fatalf("second: %v", err)
	}
	newRaw := tokenFromLink(mail.link)

	if err := svc.ResetPassword(oldRaw, "NewPass99"); !errors.Is(err, passwordreset.ErrInvalidLink) {
		t.Fatalf("old token: %v", err)
	}
	if err := svc.ResetPassword(newRaw, "NewPass99"); err != nil {
		t.Fatalf("new token: %v", err)
	}
}
