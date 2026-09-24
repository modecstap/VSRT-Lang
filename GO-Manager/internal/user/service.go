package user

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(id UserId) (*User, error) {
	return s.repo.FindByID(string(id))
}

func (s *Service) UpdateProfile(id UserId, username, email string) error {
	u, err := s.repo.FindByID(string(id))
	if err != nil {
		return err
	}
	if err := u.SetProfile(username, email); err != nil {
		return err
	}
	return s.repo.Save(u)
}

func (s *Service) SaveAvatar(id UserId, raw []byte) error {
	u, err := s.repo.FindByID(string(id))
	if err != nil {
		return err
	}
	if err := u.SetAvatar(raw); err != nil {
		return err
	}
	return s.repo.SaveAvatar(id, u.Avatar)
}
