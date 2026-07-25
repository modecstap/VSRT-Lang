package session

type Repository interface {
	Save(session *Session) (err error)
	Take(sessionId int) (session Session, err error)
}
