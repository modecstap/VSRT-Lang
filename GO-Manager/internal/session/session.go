package session

import "VSRT-Lang/internal/user"

type Session struct {
	ID      int64
	User    user.UserId
	Name    string
	Records []Record
}

type Interphase interface {
	SaveRecord(context string, window string) (record Record)
	GetRecords() []Record
}

func NewSession(
	user user.UserId,
	name string,
	translator Translator,
) *Session {
	return &Session{
		User: user,
		Name: name,
	}
}

func (s *Session) SaveRecord(context string, window string, translator Translator) (record Record) {
	record = *NewRecord(
		translator,
		window,
		context,
	)
	s.Records = append(s.Records, record)
	return record
}

func (s *Session) GetRecords() []Record {
	return s.Records
}
