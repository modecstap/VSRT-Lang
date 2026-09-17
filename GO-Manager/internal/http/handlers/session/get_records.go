package session

import (
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/session"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type recordsResponse struct {
	Records []domain.Record `json:"records"`
}

// GetRecords godoc
// @Summary      Get session records
// @Description  Retrieve records for the specified session
// @Tags         session
// @Produce      json
// @Param        id    path      int  true  "Session ID"
// @Success      200   {object}  recordsResponse
// @Failure      400   {string}  string  "invalid request"
// @Failure      404   {string}  string  "not found"
// @Security    BearerAuth
// @Router       /sessions/{id}/records [get]
func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 3 || parts[0] != "sessions" || parts[2] != "records" {
		handlers.WriteError(w, http.StatusNotFound, "not_found", "route not found")
		return
	}

	sessionID, err := strconv.Atoi(parts[1])
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid session id")
		return
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	stored, err := h.service.GetSession(userID, int64(sessionID))
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(recordsResponse{Records: stored.GetRecords()})
}
