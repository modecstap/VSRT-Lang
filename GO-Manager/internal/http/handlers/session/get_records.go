package session

import (
	"VSRT-Lang/internal/http/handlers"
	domain "VSRT-Lang/internal/session"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type recordsResponse struct {
	Records []domain.Record `json:"records"`
}

func (h *Handler) GetRecords(w http.ResponseWriter, r *http.Request) {
	// GetRecords godoc
	// @Summary      Get session records
	// @Description  Retrieve records for the specified session
	// @Tags         session
	// @Produce      json
	// @Param        id    path      int  true  "Session ID"
	// @Success      200   {object}  recordsResponse
	// @Failure      400   {string}  string  "invalid request"
	// @Failure      404   {string}  string  "not found"
	// @Router       /sessions/{id}/records [get]

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

	stored, err := h.repo.Take(sessionID)
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(recordsResponse{Records: stored.GetRecords()})
}
