package user

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
)

type MeResponse struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Avatar   *string `json:"avatar"`
}

// GetMe godoc
// @Summary      Get current user
// @Description  Return the stored username, email, and avatar for the access token subject
// @Tags         user
// @Produce      json
// @Success      200  {object}  MeResponse
// @Failure      401  {string}  string  "unauthorized"
// @Failure      404  {string}  string  "user not found"
// @Failure      500  {string}  string  "user lookup failed"
// @Security     BearerAuth
// @Router       /users/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		slog.Info("get current user", "user_id", "")
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	slog.Info("get current user", "user_id", string(userId))
	slog.Info("user service get", "user_id", string(userId))
	found, err := h.service.Get(userId)
	if err != nil {
		slog.Error("user service get failed", "user_id", string(userId), "error", err.Error())
		if err.Error() == "user not found" {
			handlers.WriteError(w, http.StatusNotFound, "user_not_found", "user not found")
			return
		}
		handlers.WriteError(w, http.StatusInternalServerError, "user_lookup_failed", err.Error())
		return
	}

	slog.Info("user service result", "username", found.Username, "email", found.Email, "avatar_present", len(found.Avatar.Bytes) > 0)
	resp := MeResponse{Username: found.Username, Email: found.Email}
	if len(found.Avatar.Bytes) > 0 {
		dataURL := "data:image/png;base64," + base64.StdEncoding.EncodeToString(found.Avatar.Bytes)
		resp.Avatar = &dataURL
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
