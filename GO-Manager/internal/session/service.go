package session

import "VSRT-Lang/internal/user"

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
	sessions, err := s.repo.TakeByUser(userId)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}