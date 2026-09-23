package session

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/session"
)

// DeleteRecord godoc
// @Summary      Delete session record
// @Description  Deletes one record from a session by phrase
// @Tags         session
// @Param        id      path  int     true  "Session ID"
// @Param        phrase  path  string  true  "Record phrase"
// @Success      204
// @Failure      400  {string}  string  "invalid request"
// @Failure      401  {string}  string  "unauthorized"
// @Failure      404  {string}  string  "not found"
// @Failure      500  {string}  string  "internal error"
// @Security     BearerAuth
// @Router       /sessions/{id}/records/{phrase} [delete]
func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 3 || len(parts) > 4 || parts[0] != "sessions" || parts[2] != "records" {
		handlers.WriteError(w, http.StatusNotFound, "not_found", "route not found")
		return
	}

	sessionID, err := strconv.Atoi(parts[1])
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid session id")
		return
	}

	phrase := ""
	if len(parts) == 4 {
		phrase = parts[3]
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	loggedUserID := ""
	if ok {
		loggedUserID = string(userID)
	}
	slog.Info("delete session record", "user_id", loggedUserID, "session_id", sessionID, "phrase", phrase)

	if strings.TrimSpace(phrase) == "" {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "phrase is required")
		return
	}
	if !ok {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	slog.Info("session service delete record", "user_id", loggedUserID, "session_id", sessionID, "phrase", phrase)
	err = h.service.DeleteRecord(userID, int64(sessionID), phrase)
	if err != nil {
		slog.Error("session service delete record failed", "user_id", loggedUserID, "session_id", sessionID, "phrase", phrase, "error", err.Error())
		if errors.Is(err, domain.ErrSessionNotFound) || errors.Is(err, domain.ErrUnauthorized) {
			handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
			return
		}
		handlers.WriteError(w, http.StatusInternalServerError, "record_delete_failed", err.Error())
		return
	}

	slog.Info("session service delete record result", "session_id", sessionID, "phrase", phrase)
	w.WriteHeader(http.StatusNoContent)
}
