package card

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	"VSRT-Lang/internal/session"
)

type cardsResponse struct {
	Cards []cardResponse `json:"cards"`
}

type cardResponse struct {
	Front cardFront `json:"front"`
	Back  cardBack  `json:"back"`
}

type cardFront struct {
	Phrase   string   `json:"phrase"`
	BaseForm string   `json:"base_form"`
	Contexts []string `json:"contexts"`
	Synonyms []string `json:"synonyms"`
}

type cardBack struct {
	Translations        []string `json:"translations"`
	ContextTranslations []string `json:"context_translations"`
}

func toCardResponse(rec session.Record) cardResponse {
	contexts := make([]string, 0, len(rec.Contexts))
	contextTranslations := make([]string, 0, len(rec.Contexts))
	for _, c := range rec.Contexts {
		contexts = append(contexts, c.Phrase)
		contextTranslations = append(contextTranslations, c.Translation)
	}

	return cardResponse{
		Front: cardFront{
			Phrase:   rec.Phrase,
			BaseForm: rec.BaseForm,
			Contexts: contexts,
			Synonyms: append(make([]string, 0, len(rec.Synonyms)), rec.Synonyms...),
		},
		Back: cardBack{
			Translations:        append(make([]string, 0, len(rec.Translations)), rec.Translations...),
			ContextTranslations: contextTranslations,
		},
	}
}

// GetCards godoc
// @Summary      Get due cards
// @Description  Returns session records that are due for review as cards
// @Tags         card
// @Produce      json
// @Param        id     path   int  true   "Session ID"
// @Param        limit  query  int  false  "Max cards"
// @Success      200  {object}  cardsResponse
// @Failure      400  {string}  string  "invalid request"
// @Failure      401  {string}  string  "unauthorized"
// @Failure      404  {string}  string  "not found"
// @Failure      500  {string}  string  "internal error"
// @Security     BearerAuth
// @Router       /sessions/{id}/cards [get]
func (h *Handler) GetCards(w http.ResponseWriter, r *http.Request) {
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

	limit := 0
	if query := r.URL.Query(); query.Has("limit") {
		limit, err = strconv.Atoi(query.Get("limit"))
		if err != nil || limit <= 0 {
			handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid limit")
			return
		}
	}

	userID, _, ok := middleware.UserFromContext(r.Context())
	loggedUserID := ""
	if ok {
		loggedUserID = string(userID)
	}
	slog.Info("get session cards", "user_id", loggedUserID, "session_id", sessionID, "limit", limit)

	if !ok {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	slog.Info("card service get cards", "user_id", loggedUserID, "session_id", sessionID, "limit", limit)
	records, err := h.service.GetCards(userID, int64(sessionID), limit)
	if err != nil {
		slog.Error("card service get cards failed", "user_id", loggedUserID, "session_id", sessionID, "limit", limit, "error", err.Error())
		if errors.Is(err, session.ErrSessionNotFound) || errors.Is(err, session.ErrUnauthorized) {
			handlers.WriteError(w, http.StatusNotFound, "session_not_found", err.Error())
			return
		}
		handlers.WriteError(w, http.StatusInternalServerError, "cards_get_failed", err.Error())
		return
	}

	cards := make([]cardResponse, 0, len(records))
	for _, rec := range records {
		cards = append(cards, toCardResponse(rec))
	}

	slog.Info("card service get cards result", "session_id", sessionID, "count", len(cards))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(cardsResponse{Cards: cards})
}
