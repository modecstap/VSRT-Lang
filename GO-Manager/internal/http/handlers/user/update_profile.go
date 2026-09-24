package user

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/user"
)

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UpdateProfile godoc
// @Summary      Update current user profile
// @Description  Set username and email for the access token subject. Codes: username_required, email_required, invalid_email, invalid_request, unauthorized, username_taken, email_taken, user_not_found, profile_save_failed
// @Tags         user
// @Accept       json
// @Param        body  body  UpdateProfileRequest  true  "Username and email"
// @Success      204
// @Failure      400  {string}  string  "invalid profile"
// @Failure      401  {string}  string  "unauthorized"
// @Failure      409  {string}  string  "username or email taken"
// @Security     BearerAuth
// @Router       /users/me [post]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		slog.Info("update profile", "user_id", "")
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	slog.Info("update profile", "user_id", string(userId), "username", req.Username, "email", req.Email)
	err := h.service.UpdateProfile(userId, req.Username, req.Email)
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	switch {
	case errors.Is(err, domain.ErrUsernameRequired):
		handlers.WriteError(w, http.StatusBadRequest, "username_required", "username required")
	case errors.Is(err, domain.ErrEmailRequired):
		handlers.WriteError(w, http.StatusBadRequest, "email_required", "email required")
	case errors.Is(err, domain.ErrInvalidEmail):
		handlers.WriteError(w, http.StatusBadRequest, "invalid_email", "invalid email")
	case errors.Is(err, domain.ErrUsernameTaken):
		handlers.WriteError(w, http.StatusConflict, "username_taken", "username taken")
	case errors.Is(err, domain.ErrEmailTaken):
		handlers.WriteError(w, http.StatusConflict, "email_taken", "email taken")
	case errors.Is(err, domain.ErrNotFound):
		handlers.WriteError(w, http.StatusNotFound, "user_not_found", "user not found")
	default:
		handlers.WriteError(w, http.StatusInternalServerError, "profile_save_failed", err.Error())
	}
}
