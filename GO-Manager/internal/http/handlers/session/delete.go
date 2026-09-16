package session

import (
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/user"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

// DeleteSession godoc
// @Summary     Delete session
// @Description Deletes a session by its ID.
// @Tags        session
// @Param       id path int true "Session ID"
// @Success     204
// @Failure     400 {string} string "Invalid session ID"
// @Failure     404 {string} string "Session not found"
// @Failure     500 {string} string "Internal server error"
// @Security    BearerAuth
// @Router      /sessions/{id} [delete]
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "sessions" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	sessionID, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	sessionsId, err := h.takeSessionsIdByUser(userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if !slices.Contains(sessionsId, sessionID) {
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	err = h.repo.Delete(sessionID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) takeSessionsIdByUser(userID user.UserId) ([]int, error) {
	sessions, err := h.repo.FindByUser(userID)
	if err != nil {
		return nil, err
	}
	sessionsId := make([]int, len(sessions))
	for i, session := range sessions {
		sessionsId[i] = int(session.ID)
	}
	return sessionsId, nil
}
