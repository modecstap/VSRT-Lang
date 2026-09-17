package memory

import (
	"VSRT-Lang/internal/session"
	"VSRT-Lang/internal/user"
	"database/sql"
	"sort"
	"sync"
)

type SessionRepository struct {
	mu          sync.RWMutex
	sessions    map[int64]*session.Session
	nextSession int64
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		sessions:    make(map[int64]*session.Session),
		nextSession: 1,
	}
}

func (r *SessionRepository) Save(s *session.Session) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if s.ID == 0 {
		s.ID = r.nextSession
		r.nextSession++
	} else if s.ID >= r.nextSession {
		r.nextSession = s.ID + 1
	}

	r.sessions[s.ID] = cloneSession(s)
	return s.ID, nil
}

func (r *SessionRepository) Take(sessionID int64) (session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	s, ok := r.sessions[sessionID]
	if !ok {
		return session.Session{}, sql.ErrNoRows
	}
	return *cloneSession(s), nil
}

func (r *SessionRepository) FindByUser(userID user.UserId) ([]session.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]int64, 0)
	for id, s := range r.sessions {
		if s.User == userID {
			ids = append(ids, id)
		}
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	result := make([]session.Session, 0, len(ids))
	for _, id := range ids {
		result = append(result, *cloneSession(r.sessions[id]))
	}
	return result, nil
}

func (r *SessionRepository) Delete(sessionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sessions, sessionID)
	return nil
}

func cloneSession(src *session.Session) *session.Session {
	dst := &session.Session{
		ID:      src.ID,
		User:    src.User,
		Name:    src.Name,
		Records: make(map[string]*session.Record, len(src.Records)),
	}
	for phrase, record := range src.Records {
		dst.Records[phrase] = cloneRecord(record)
	}
	return dst
}

func cloneRecord(src *session.Record) *session.Record {
	if src == nil {
		return nil
	}
	dst := *src
	dst.Translations = append([]string(nil), src.Translations...)
	dst.Synonyms = append([]string(nil), src.Synonyms...)
	dst.Antonyms = append([]string(nil), src.Antonyms...)
	dst.Contexts = append([]session.Context(nil), src.Contexts...)
	return &dst
}
