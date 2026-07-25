package session

type UserId int

type Session struct {
	ID      int64
	User    UserId
	Name    string
	Records []Record

	Translator Translator
}

type Interphase interface {
	SaveRecord(context string, window string) (record Record)
}

func (s *Session) SaveRecord(context string, window string) (record Record) {
	record = *NewRecord(
		s.Translator,
		window,
		context,
	)
	s.Records = append(s.Records, record)
	return record
}
