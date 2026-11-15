package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (h *Handler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	listID, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}
	userID, err := parseUserIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))

	resp, err := h.permissions.GetUserPermissions(ctx, userID, listID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			WriteErrorJSON(w, err, 401)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			WriteErrorJSON(w, err, 403)
			return
		}
		if errors.Is(err, domain.ErrUserPermissionsNotFound) || errors.Is(err, domain.ErrUserNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to get user permissions", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
