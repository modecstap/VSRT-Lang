package memory

import (
	"sort"
	"sync"

	"VSRT-Lang/internal/card"
)

type KnowledgeRepository struct {
	mu        sync.RWMutex
	knowledge map[card.KnowledgeId]card.Knowledge
}

func NewKnowledgeRepository() *KnowledgeRepository {
	return &KnowledgeRepository{
		knowledge: make(map[card.KnowledgeId]card.Knowledge),
	}
}

func (r *KnowledgeRepository) Save(k *card.Knowledge) (card.KnowledgeId, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.knowledge[k.ID] = *k
	return k.ID, nil
}

func (r *KnowledgeRepository) Take(id card.KnowledgeId) (*card.Knowledge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	stored, ok := r.knowledge[id]
	if !ok {
		return nil, nil
	}
	return &stored, nil
}

func (r *KnowledgeRepository) FindBySession(sessionID int64) ([]card.Knowledge, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]card.Knowledge, 0)
	for id, k := range r.knowledge {
		if id.SessionID == sessionID {
			result = append(result, k)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID.PhraseKey < result[j].ID.PhraseKey })
	return result, nil
}
