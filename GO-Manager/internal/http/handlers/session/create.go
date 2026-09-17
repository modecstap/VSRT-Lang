package session

import (
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	"encoding/json"
	"net/http"
)

type sessionResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type createSessionRequest struct {
	Name string `json:"name"`
}

// CreateSession godoc
//
// @Summary     Create session
// @Description Create a new session for the authenticated user
// @Tags        session
// @Accept      json
// @Produce     json
// @Param       body body createSessionRequest true "Session data"
// @Success     201 {object} sessionResponse
// @Failure     400 {string} string "invalid request"
// @Failure     401 {string} string "unauthorized"
// @Security    BearerAuth
// @Router      /sessions [post]
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var req createSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	if !ok || userID == "" {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	id, err := h.service.NewSession(userID, req.Name)
	if err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, "session_save_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sessionResponse{ID: id, Name: req.Name})
}
