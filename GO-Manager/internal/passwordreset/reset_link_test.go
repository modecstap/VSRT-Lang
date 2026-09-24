package passwordreset

import (
	"testing"
	"time"

	"VSRT-Lang/internal/user"
)

func TestResetLink_ActiveUntilTTL(t *testing.T) {
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	link := NewResetLink(user.UserId("u1"), "hash-a", now)

	if !link.Active(now) {
		t.Fatal("new link must be active")
	}
	if !link.Active(now.Add(LinkTTL - time.Second)) {
		t.Fatal("link must stay active before TTL")
	}
	if link.Active(now.Add(LinkTTL)) {
		t.Fatal("link must be inactive at TTL boundary")
	}
}

func TestResetLink_ConsumeKeepsCooldown(t *testing.T) {
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	link := NewResetLink(user.UserId("u1"), "hash-a", now)
	sentAt := link.SentAt

	link.Consume()
	if link.TokenHash != "" || link.ExpiresAt != nil {
		t.Fatal("Consume must clear hash and expiry")
	}
	if !link.SentAt.Equal(sentAt) {
		t.Fatal("Consume must keep SentAt")
	}
	if link.Active(now) {
		t.Fatal("consumed link must not be active")
	}
	if !link.CoolingDown(now.Add(time.Minute)) {
		t.Fatal("cooldown must remain after Consume")
	}
}

func TestResetLink_CooldownBoundary(t *testing.T) {
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	link := NewResetLink(user.UserId("u1"), "hash-a", now)

	if !link.CoolingDown(now.Add(SendCooldown - time.Second)) {
		t.Fatal("must cool down before boundary")
	}
	if link.CoolingDown(now.Add(SendCooldown)) {
		t.Fatal("cooldown must be false at boundary")
	}
}

func TestResetLink_DifferentHashesAreIndependent(t *testing.T) {
	now := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	old := NewResetLink(user.UserId("u1"), "hash-old", now)
	newer := NewResetLink(user.UserId("u1"), "hash-new", now.Add(time.Minute))

	old.Consume()
	if old.Active(now.Add(time.Minute)) {
		t.Fatal("old consumed hash must stay dead")
	}
	if !newer.Active(now.Add(time.Minute)) {
		t.Fatal("new hash entity stays active on its own")
	}
}
