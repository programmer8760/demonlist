package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type CreatePermissionsRequest struct {
	UserID      uint   `json:"user_id"`
	Permissions uint16 `json:"permissions"`
}

func (h *Handler) CreatePermissions(w http.ResponseWriter, r *http.Request) {
	listID, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	var req CreatePermissionsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.permissions.AddUserPermissions(ctx, domain.UserPermissions{
		UserID:      req.UserID,
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
		if errors.Is(err, domain.ErrListNotFound) || errors.Is(err, domain.ErrUserNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		if errors.Is(err, domain.ErrUserPermissionsAlreadyExist) {
			WriteErrorJSON(w, err, 409)
			return
		}
		log.Error("failed to create user permissions", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
