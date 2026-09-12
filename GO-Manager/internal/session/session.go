package session

import "VSRT-Lang/internal/user"

type Session struct {
	ID      int64
	User    user.UserId
	Name    string
	Records map[string]Record
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
		Records: make(map[string]Record),
	}
}

func (s *Session) SaveRecord(context string, window string, translator Translator) (record Record) {
	record, ok := s.Records[window]
	if !ok {
		record = *NewRecord(
			translator,
			window,
			context,
		)
	}
	s.Records[window] = record
	return record
}

func (s *Session) GetRecords() []Record {
	records := make([]Record, 0, len(s.Records))
	for _, record := range s.Records {
		records = append(records, record)
	}
	return records
}
