package user

import (
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	_ "VSRT-Lang/internal/session"
	"encoding/json"
	"net/http"
	"strings"
)

// GetUserSessions godoc
// @Summary      Get user sessions
// @Description  Retrieve all sessions for the authenticated user
// @Tags         user
// @Produce      json
// @Success      200      {array}   session.Session
// @Failure      404      {string}  string  "not found"
// @Failure      500      {string}  string  "internal server error"
// @Security    BearerAuth
// @Router       /users/{user_id}/sessions [get]
func (h *Handler) GetUserSessions(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 || parts[0] != "users" || parts[1] != "sessions" {
		handlers.WriteError(w, http.StatusNotFound, "not_found", "route not found")
		return
	}

	userId, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		handlers.WriteError(w, http.StatusInternalServerError, "invalid_token", "")
		return
	}

	sessions, err := h.sessionRepo.FindByUser(userId)
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, "sessions_not_found", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sessions)
}
