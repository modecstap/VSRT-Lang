package card

import (
	"time"

	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
)

type Service struct {
	sessions  session.Repository
	knowledge Repository
	now       func() time.Time
}

func NewService(sessions session.Repository, knowledge Repository, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{
		sessions:  sessions,
		knowledge: knowledge,
		now:       now,
	}
}

func (s *Service) GetCards(userID user.UserId, sessionID int64, limit int) ([]session.Record, error) {
	sess, err := s.findSession(userID, sessionID)
	if err != nil {
		return nil, err
	}

	known, err := s.knowledge.FindBySession(sessionID)
	if err != nil {
		return nil, err
	}

	return DueRecords(sess.GetRecords(), known, s.now().UTC(), limit), nil
}

type UpdateKnowledgeCommand struct {
	UserId     user.UserId
	SessionId  int64
	Phrase     string
	Estimation Estimation
}

func (s *Service) UpdateKnowledge(c UpdateKnowledgeCommand) error {
	sess, err := s.findSession(c.UserId, c.SessionId)
	if err != nil {
		return err
	}

	key := session.RecordKey(c.Phrase)
	if _, ok := sess.Records[key]; !ok {
		return ErrCardNotFound
	}

	id := KnowledgeId{SessionID: c.SessionId, PhraseKey: key}
	k, err := s.knowledge.Take(id)
	if err != nil {
		return err
	}
	if k == nil {
		k = NewKnowledge(id)
	}

	if err := k.Estimate(c.Estimation, s.now().UTC()); err != nil {
		return err
	}

	_, err = s.knowledge.Save(k)
	return err
}

func (s *Service) findSession(userID user.UserId, sessionID int64) (*session.Session, error) {
	sessions, err := s.sessions.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	for i := range sessions {
		if sessions[i].ID == sessionID {
			return &sessions[i], nil
		}
	}

	return nil, session.ErrSessionNotFound
}
