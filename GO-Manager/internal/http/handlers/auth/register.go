package auth

import (
	"encoding/json"
	"net/http"
)

// Register godoc
// @Summary      Register user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "Registration data"
// @Success      200   {object}  map[string]any
// @Failure      400   {string}  string  "invalid request or validation error"
// @Router       /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{"id": user.ID, "email": user.Email})
}
