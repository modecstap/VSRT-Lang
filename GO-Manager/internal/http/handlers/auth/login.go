package auth

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tokens, err := h.service.Login(req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	_ = json.NewEncoder(w).Encode(tokens)
}
