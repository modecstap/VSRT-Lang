package user

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"VSRT-Lang/internal/http/handlers"
	"VSRT-Lang/internal/http/middleware"
	domain "VSRT-Lang/internal/user"
)

const multipartEnvelopeSlack = 32 * 1024

// SaveAvatar godoc
// @Summary      Save user avatar
// @Description  Upload a JPEG, PNG, or WebP image; it is cropped to a square if needed and stored as PNG
// @Tags         user
// @Accept       mpfd
// @Param        avatar  formData  file  true  "Avatar image"
// @Success      204
// @Failure      400      {string}  string  "invalid avatar"
// @Failure      401      {string}  string  "unauthorized"
// @Security     BearerAuth
// @Router       /users/avatar [post]
func (h *Handler) SaveAvatar(w http.ResponseWriter, r *http.Request) {
	userId, _, ok := middleware.UserFromContext(r.Context())
	if !ok {
		handlers.WriteError(w, http.StatusUnauthorized, "unauthorized", "missing user context")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, domain.MaxAvatarBytes+multipartEnvelopeSlack)
	if err := r.ParseMultipartForm(domain.MaxAvatarBytes); err != nil {
		if isAvatarTooLarge(err) {
			handlers.WriteError(w, http.StatusBadRequest, "avatar_too_large", err.Error())
			return
		}
		handlers.WriteError(w, http.StatusBadRequest, "avatar_missing", err.Error())
		return
	}

	file, _, err := r.FormFile("avatar")
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "avatar_missing", "avatar file required")
		return
	}
	defer file.Close()

	raw, err := io.ReadAll(file)
	if err != nil {
		if isAvatarTooLarge(err) {
			handlers.WriteError(w, http.StatusBadRequest, "avatar_too_large", err.Error())
			return
		}
		handlers.WriteError(w, http.StatusBadRequest, "avatar_missing", err.Error())
		return
	}

	if err := h.service.SaveAvatar(userId, raw); err != nil {
		switch {
		case errors.Is(err, domain.ErrAvatarTooLarge):
			handlers.WriteError(w, http.StatusBadRequest, "avatar_too_large", err.Error())
		case errors.Is(err, domain.ErrAvatarDimensions):
			handlers.WriteError(w, http.StatusBadRequest, "avatar_dimensions_invalid", err.Error())
		case errors.Is(err, domain.ErrInvalidAvatarType):
			handlers.WriteError(w, http.StatusBadRequest, "invalid_avatar_type", err.Error())
		case err.Error() == "user not found":
			handlers.WriteError(w, http.StatusNotFound, "user_not_found", err.Error())
		default:
			handlers.WriteError(w, http.StatusInternalServerError, "avatar_save_failed", err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isAvatarTooLarge(err error) bool {
	var maxBytes *http.MaxBytesError
	if errors.As(err, &maxBytes) {
		return true
	}
	return errors.Is(err, multipart.ErrMessageTooLarge)
}
