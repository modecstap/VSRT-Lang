package passwordreset

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"
	"time"

	"VSRT-Lang/internal/user"
)

var (
	ErrRateLimited = errors.New("try again in 15 minutes")
	ErrSendFailed  = errors.New("Unable to send email")
	ErrInvalidLink = errors.New("This link is invalid or expired")
)

type Service struct {
	users         user.Repository
	links         Repository
	mailer        Mailer
	tx            Transactor
	publicSiteURL string
	now           func() time.Time
}

func NewService(
	users user.Repository,
	links Repository,
	mailer Mailer,
	tx Transactor,
	publicSiteURL string,
	now func() time.Time,
) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{
		users:         users,
		links:         links,
		mailer:        mailer,
		tx:            tx,
		publicSiteURL: publicSiteURL,
		now:           now,
	}
}

func (s *Service) RequestReset(email string) error {
	found, err := s.users.FindByEmail(email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil
		}
		return err
	}

	existing, err := s.links.Take(user.UserId(found.ID))
	if err != nil {
		return err
	}
	now := s.now()
	if existing != nil && existing.CoolingDown(now) {
		return ErrRateLimited
	}

	raw, err := newRawToken()
	if err != nil {
		return err
	}
	hash := hashToken(raw)

	linkURL, err := s.buildResetURL(raw)
	if err != nil {
		return ErrSendFailed
	}

	if err := s.mailer.Send(found.Email, linkURL); err != nil {
		return ErrSendFailed
	}

	sentAt := s.now()
	newLink := NewResetLink(user.UserId(found.ID), hash, sentAt)

	return s.tx.Within(func(repos TxRepos) error {
		current, err := repos.Links.TakeForUpdate(user.UserId(found.ID))
		if err != nil {
			return err
		}
		if current != nil && current.CoolingDown(s.now()) {
			return ErrRateLimited
		}
		_, err = repos.Links.Save(&newLink)
		return err
	})
}

func (s *Service) ResetPassword(rawToken, plain string) error {
	return s.tx.Within(func(repos TxRepos) error {
		if rawToken == "" {
			return ErrInvalidLink
		}

		hash := hashToken(rawToken)
		link, err := repos.Links.TakeByHash(hash)
		if err != nil {
			return err
		}
		now := s.now()
		if link == nil || !link.Active(now) {
			return ErrInvalidLink
		}

		found, err := repos.Users.FindByID(string(link.User))
		if err != nil {
			if errors.Is(err, user.ErrNotFound) {
				return ErrInvalidLink
			}
			return err
		}

		if err := found.SetPassword(plain); err != nil {
			return err
		}

		link.Consume()
		if err := repos.Users.Save(found); err != nil {
			return err
		}
		if _, err := repos.Links.Save(link); err != nil {
			return err
		}
		return repos.Refresh.RevokeByUserID(found.ID)
	})
}

func (s *Service) buildResetURL(rawToken string) (string, error) {
	base := strings.TrimRight(s.publicSiteURL, "/")
	parsed, err := url.Parse(base)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", ErrSendFailed
	}
	q := url.Values{}
	q.Set("token", rawToken)
	return base + "/reset-password?" + q.Encode(), nil
}

func newRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
