package session

import (
	"VSRT-Lang/internal/user"
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized")
var ErrSessionNotFound = errors.New("session not found")

type Service struct {
	repo       Repository
	translator Translator
}

func (s *Service) DeleteRecord(userId user.UserId, sessionId int64, phrase string) error {
	sessions, err := s.repo.TakeByUser(userId)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.ID == sessionId {
			delete(session.Records, phrase)
			_, err = s.repo.Save(&session)
			return err
		}
	}

	return ErrUnauthorized
}

func NewService(repo Repository, translator Translator) *Service {
	return &Service{
		repo:       repo,
		translator: translator,
	}
}

func (s *Service) NewSession(userId user.UserId, name string) (int64, error) {
	session := NewSession(userId, name)
	id, err := s.repo.Save(session)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) GetSessions(userId user.UserId) ([]Session, error) {
	sessions, err := s.repo.TakeByUser(userId)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

type AddRecordCommand struct {
	UserId    user.UserId
	SessionId int64
	Phrase    string
	Context   string
}

func (s *Service) AddRecord(c AddRecordCommand) (Record, error) {
	sessions, err := s.repo.TakeByUser(c.UserId)
	if err != nil {
		return Record{}, err
	}

	for _, session := range sessions {
		if session.ID == c.SessionId {
			record := session.SaveRecord(c.Context, c.Phrase, s.translator)
			s.repo.Save(&session)
			return record, nil
		}
	}

	return Record{}, ErrSessionNotFound
}

func (s *Service) Delete(userId user.UserId, sessionId int64) error {
	sessions, err := s.repo.TakeByUser(userId)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.ID == sessionId {
			return s.repo.Delete(sessionId)
		}
	}

	return ErrUnauthorized
}
