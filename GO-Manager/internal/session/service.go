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
	sessions, err := s.repo.FindByUser(userId)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (s *Service) GetSession(userId user.UserId, sessionId int64) (Session, error) {
	session, err := s.findSession(userId, sessionId)
	if err != nil {
		return Session{}, err
	}

	return *session, nil
}

func (s *Service) DeleteSession(userId user.UserId, sessionId int64) error {
	session, err := s.findSession(userId, sessionId)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ErrUnauthorized
		}
		return err
	}

	return s.repo.Delete(session.ID)
}

type AddRecordCommand struct {
	UserId    user.UserId
	SessionId int64
	Phrase    string
	Context   string
}

func (s *Service) AddRecord(c AddRecordCommand) (Record, error) {
	session, err := s.findSession(c.UserId, c.SessionId)
	if err != nil {
		return Record{}, err
	}

	record, err := session.SaveRecord(c.Context, c.Phrase, s.translator)
	if err != nil {
		return Record{}, err
	}

	if _, err := s.repo.Save(session); err != nil {
		return Record{}, err
	}

	return record, nil
}

func (s *Service) DeleteRecord(userId user.UserId, sessionId int64, phrase string) error {
	session, err := s.findSession(userId, sessionId)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return ErrUnauthorized
		}
		return err
	}

	delete(session.Records, phrase)
	_, err = s.repo.Save(session)
	return err
}

func (s *Service) findSession(userID user.UserId, sessionID int64) (*Session, error) {
	sessions, err := s.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	for i := range sessions {
		if sessions[i].ID == sessionID {
			return &sessions[i], nil
		}
	}

	return nil, ErrSessionNotFound
}
