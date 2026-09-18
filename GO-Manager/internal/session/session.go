package session

import "VSRT-Lang/internal/user"

type Session struct {
	ID      int64
	User    user.UserId
	Name    string
	Records map[string]*Record
}

type Interphase interface {
	SaveRecord(context string, window string) (record Record)
	GetRecords() []Record
}

func NewSession(
	user user.UserId,
	name string,
) *Session {
	return &Session{
		User:    user,
		Name:    name,
		Records: make(map[string]*Record),
	}
}

func (s *Session) SaveRecord(context string, window string, translator Translator) (Record, error) {
	recordPtr, ok := s.Records[window]
	if !ok {
		newRecord, err := NewRecord(translator, window, context)
		if err != nil {
			return Record{}, err
		}
		s.Records[window] = newRecord
		return *s.Records[window], nil
	}

	recordPtr.Count++

	return *recordPtr, nil
}

func (s *Session) GetRecords() []Record {
	records := make([]Record, 0, len(s.Records))
	for _, record := range s.Records {
		records = append(records, *record)
	}
	return records
}
