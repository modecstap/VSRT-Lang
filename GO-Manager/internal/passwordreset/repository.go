package passwordreset

import "VSRT-Lang/internal/user"

type Repository interface {
	Save(link *ResetLink) (user.UserId, error)
	Take(id user.UserId) (*ResetLink, error)
	TakeForUpdate(id user.UserId) (*ResetLink, error)
	TakeByHash(hash string) (*ResetLink, error)
}

type TxRepos struct {
	Users   user.Repository
	Links   Repository
	Refresh interface{ RevokeByUserID(userID string) error }
}

type Transactor interface {
	Within(fn func(TxRepos) error) error
}

type Mailer interface {
	Send(to, link string) error
}
