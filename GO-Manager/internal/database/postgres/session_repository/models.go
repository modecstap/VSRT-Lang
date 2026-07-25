package session_repository

import "VSRT-Lang/internal/session"

type dbContext struct {
	Phrase      string `json:"phrase"`
	Translation string `json:"translation"`
}

func toDBContexts(src []session.Context) []dbContext {
	dst := make([]dbContext, len(src))
	for i, c := range src {
		dst[i] = dbContext(c)
	}
	return dst
}

func fromDBContexts(src []dbContext) []session.Context {
	dst := make([]session.Context, len(src))
	for i, c := range src {
		dst[i] = session.Context(c)
	}
	return dst
}
