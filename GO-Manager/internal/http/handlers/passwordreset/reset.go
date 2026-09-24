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

// ResetPassword godoc
// @Summary      Reset password with link token
// @Description  Sets a new password when the token is valid and not expired.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  ResetPasswordRequest  true  "Token and new password"
// @Success      204
// @Failure      400  {object}  handlers.ErrorResponse
// @Failure      500  {object}  handlers.ErrorResponse
// @Router       /reset-password [post]
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	slog.Info("reset password input", "token", "[REDACTED]", "password", "[REDACTED]")

	if req.Token == "" {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_link", passwordreset.ErrInvalidLink.Error())
		return
	}
	if !user.ValidatePassword(req.Password) {
		handlers.WriteError(w, http.StatusBadRequest, "weak_password", user.ErrWeakPassword.Error())
		return
	}

	slog.Info("reset password command", "token", "[REDACTED]")
	err := h.service.ResetPassword(req.Token, req.Password)
	if err == nil {
		slog.Info("reset password result", "result", "ok")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	slog.Info("reset password result", "result", err.Error())
	switch {
	case errors.Is(err, passwordreset.ErrInvalidLink):
		handlers.WriteError(w, http.StatusBadRequest, "invalid_link", err.Error())
	case errors.Is(err, user.ErrWeakPassword):
		handlers.WriteError(w, http.StatusBadRequest, "weak_password", err.Error())
	default:
		slog.Error("reset password internal", "error", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal", "invalid request")
	}
}
