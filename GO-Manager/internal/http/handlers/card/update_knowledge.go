package card

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	domain "VSRT-Lang/internal/card"
	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/session"
)

type estimateCardRequest struct {
	Phrase     string `json:"phrase"`
	Estimation string `json:"estimation"`
}

// UpdateKnowledge godoc
// @Summary      Estimate card
// @Description  Applies a review estimation (REPEAT, DIFFICULT, EASY, MOMENTAL) to a session card
// @Tags         card
// @Accept       json
// @Param        id    path  int                  true  "Session ID"
// @Param        body  body  estimateCardRequest  true  "Card estimation"
// @Success      204
// @Failure      400  {string}  string  "invalid request"
// @Failure      401  {string}  string  "unauthorized"
// @Failure      404  {string}  string  "not found"
// @Failure      500  {string}  string  "internal error"
// @Security     BearerAuth
// @Router       /sessions/{id}/cards [post]
func (h *Handler) UpdateKnowledge(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "sessions" || parts[2] != "cards" {
		handlers.WriteError(w, http.StatusNotFound, "not_found", "route not found")
		return
	}

	sessionID, err := strconv.Atoi(parts[1])
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid session id")
		return
	}

	var req estimateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	loggedUserID := ""
	if ok {
		loggedUserID = string(userID)
	}
	slog.Info("estimate session card", "user_id", loggedUserID, "session_id", sessionID, "phrase", req.Phrase, "estimation", req.Estimation)

	if strings.TrimSpace(req.Phrase) == "" {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "phrase is required")
		return
	}
	if !ok {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	slog.Info("card service update knowledge", "user_id", loggedUserID, "session_id", sessionID, "phrase", req.Phrase, "estimation", req.Estimation)
	err = h.service.UpdateKnowledge(domain.UpdateKnowledgeCommand{
		UserId:     userID,
		SessionId:  int64(sessionID),
		Phrase:     req.Phrase,
		Estimation: domain.Estimation(req.Estimation),
	})
	if err != nil {
		slog.Error("card service update knowledge failed", "user_id", loggedUserID, "session_id", sessionID, "phrase", req.Phrase, "estimation", req.Estimation, "error", err.Error())
		switch {
		case errors.Is(err, session.ErrSessionNotFound), errors.Is(err, session.ErrUnauthorized):
			handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
		case errors.Is(err, domain.ErrCardNotFound):
			handlers.WriteError(w, http.StatusNotFound, "card_not_found", err.Error())
		case errors.Is(err, domain.ErrUnknownEstimation):
			handlers.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		default:
			handlers.WriteError(w, http.StatusInternalServerError, "card_estimate_failed", err.Error())
		}
		return
	}

	slog.Info("card service update knowledge result", "session_id", sessionID, "phrase", req.Phrase, "estimation", req.Estimation)
	w.WriteHeader(http.StatusNoContent)
}
