package passwordreset

import (
	"time"

	"VSRT-Lang/internal/user"
)

const (
	LinkTTL      = 15 * time.Minute
	SendCooldown = 15 * time.Minute
)

type ResetLink struct {
	User      user.UserId
	TokenHash string
	SentAt    time.Time
	ExpiresAt *time.Time
}

func NewResetLink(userID user.UserId, tokenHash string, now time.Time) ResetLink {
	expires := now.Add(LinkTTL)
	return ResetLink{
		User:      userID,
		TokenHash: tokenHash,
		SentAt:    now,
		ExpiresAt: &expires,
	}
}

func (l ResetLink) Active(now time.Time) bool {
	if l.TokenHash == "" || l.ExpiresAt == nil {
		return false
	}
	return now.Before(*l.ExpiresAt)
}

func (l ResetLink) CoolingDown(now time.Time) bool {
	if l.SentAt.IsZero() {
		return false
	}
	return now.Before(l.SentAt.Add(SendCooldown))
}

func (l *ResetLink) Consume() {
	l.TokenHash = ""
	l.ExpiresAt = nil
}
