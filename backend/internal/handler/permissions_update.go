package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type UpdatePermissionsRequest struct {
	Permissions uint16 `json:"permissions"`
}

func (h *Handler) UpdatePermissions(w http.ResponseWriter, r *http.Request) {
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

	var req UpdatePermissionsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.permissions.UpdateUserPermissions(ctx, domain.UserPermissions{
		UserID:      userID,
		ListID:      listID,
		Permissions: req.Permissions,
	})
	if err != nil {
		var br domain.BadRequest
		if errors.As(err, &br) {
			WriteErrorJSON(w, err, 400)
			return
		}
		if errors.Is(err, domain.ErrUnauthorized) {
			WriteErrorJSON(w, err, 401)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			WriteErrorJSON(w, err, 403)
			return
		}
		if errors.Is(err, domain.ErrUserPermissionsNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to update user permissions", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
