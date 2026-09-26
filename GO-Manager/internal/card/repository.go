package card

type Repository interface {
	Save(k *Knowledge) (KnowledgeId, error)
	// Take returns (nil, nil) when no knowledge is stored for id.
	Take(id KnowledgeId) (*Knowledge, error)
	FindBySession(sessionID int64) ([]Knowledge, error)
}
