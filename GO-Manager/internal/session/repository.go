package session

type Repository interface {
	Save(session *Session) (id int64, err error)
	Take(sessionId int64) (session Session, err error)
}
