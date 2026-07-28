package auth

import (
	"encoding/json"
	"net/http"

	handlers "VSRT-Lang/internal/http/handlers"
)

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "Login credentials"
// @Success      200   {object}  map[string]any
// @Failure      400   {string}  string  "invalid request"
// @Failure      401   {string}  string  "unauthorized"
// @Router       /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid_request", "invalid request")
		return
	}

	tokens, err := h.service.Login(req.Login, req.Password)
	if err != nil {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(tokens)
}
