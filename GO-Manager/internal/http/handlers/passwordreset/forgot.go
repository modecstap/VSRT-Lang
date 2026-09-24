package passwordreset

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	handlers "VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/passwordreset"
	"VSRT-Lang/internal/user"
)

// ForgotPassword godoc
// @Summary      Request password reset email
// @Description  Sends a reset link when the email belongs to an account. Always returns 204 for unknown emails.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  ForgotPasswordRequest  true  "Account email"
// @Success      204
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /forgot-password [post]
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	slog.Info("forgot password input", "email", req.Email)

	if !user.ValidateEmail(req.Email) {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_email", "invalid email")
		return
	}

	slog.Info("forgot password command", "email", req.Email)
	err := h.service.RequestReset(req.Email)
	if err == nil {
		slog.Info("forgot password result", "result", "ok")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	slog.Info("forgot password result", "result", err.Error())
	switch {
	case errors.Is(err, passwordreset.ErrRateLimited):
		handlers.WriteError(w, http.StatusBadRequest, "rate_limited", err.Error())
	case errors.Is(err, passwordreset.ErrSendFailed):
		handlers.WriteError(w, http.StatusBadRequest, "email_failed", err.Error())
	default:
		slog.Error("forgot password internal", "error", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal", "invalid request")
	}
}
