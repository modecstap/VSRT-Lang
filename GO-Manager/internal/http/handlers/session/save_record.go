package session

import (
	"VSRT-Lang/internal"
	"VSRT-Lang/internal/http/handlers"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type createRecordRequest struct {
	Phrase  string `json:"phrase"`
	Context string `json:"context"`
}

func (h *Handler) SaveRecord(w http.ResponseWriter, r *http.Request) {
	// SaveRecord godoc
	// @Summary      Save record to session
	// @Description  Save a new record into the specified session
	// @Tags         session
	// @Accept       json
	// @Produce      json
	// @Param        id    path      int                 true  "Session ID"
	// @Param        body  body      createRecordRequest true  "Record data"
	// @Success      201   {object}  map[string]any
	// @Failure      400   {string}  string  "invalid request"
	// @Failure      404   {string}  string  "not found"
	// @Failure      500   {string}  string  "internal error"
	// @Router       /sessions/{id}/records [post]

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

	var req createRecordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	stored, err := h.repo.Take(sessionID)
	if err != nil {
		handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
		return
	}

	stored.Translator = internal.MockTranslator{}
	stored.SaveRecord(req.Phrase, req.Context)
	if err := h.repo.Save(&stored); err != nil {
		handlers.WriteError(w, http.StatusInternalServerError, "record_save_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": "record saved"})
}
