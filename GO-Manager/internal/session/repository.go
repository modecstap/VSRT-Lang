package session

import "VSRT-Lang/internal/user"

type Repository interface {
	Save(session *Session) (id int64, err error)
	Take(sessionId int64) (session Session, err error)
	TakeByUser(userId user.UserId) ([]Session, error)
	Delete(sessionId int64) error
}
