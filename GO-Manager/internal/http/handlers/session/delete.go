package session

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/session"
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
	err = h.service.DeleteSession(userID, int64(sessionID))
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) || errors.Is(err, domain.ErrSessionNotFound) {
			handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
			return
		}
		handlers.WriteError(w, http.StatusInternalServerError, "session_delete_failed", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
