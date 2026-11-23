package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/service"
)

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUserIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	var req UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))

	resp, err := h.users.UpdateUser(ctx, service.UpdateUserInput{
		ID:       id,
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
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
		if errors.Is(err, domain.ErrUserNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			WriteErrorJSON(w, err, 409)
			return
		}
		log.Error("failed to update user", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
