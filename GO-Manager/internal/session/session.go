package session

import (
	"strings"

	"VSRT-Lang/internal/user"
)

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

func (s *Session) SaveRecord(context string, phrase string, translator Translator) (Record, error) {
	key := recordKey(phrase)
	if existing, ok := s.Records[key]; ok {
		if err := existing.registerSave(context, translator); err != nil {
			return Record{}, err
		}
		return *existing, nil
	}

	newRecord, err := NewRecord(
		NewRecordCommand{
			Translator:  translator,
			Phrase:      phrase,
			MainContext: context,
		},
	)
	if err != nil {
		return Record{}, err
	}
	s.Records[key] = newRecord
	return *newRecord, nil
}

func (s *Session) DeleteRecord(phrase string) {
	delete(s.Records, recordKey(phrase))
}

func recordKey(phrase string) string {
	return strings.ToLower(strings.TrimSpace(phrase))
}

func (s *Session) GetRecords() []Record {
	records := make([]Record, 0, len(s.Records))
	for _, record := range s.Records {
		records = append(records, *record)
	}
	return records
}
