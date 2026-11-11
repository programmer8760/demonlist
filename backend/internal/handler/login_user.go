package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.users.LoginUser(ctx, domain.User{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		var br domain.BadRequest
		if errors.As(err, &br) {
			WriteErrorJSON(w, err, 400)
			return
		}
		if errors.Is(err, domain.ErrInvalidPassword) {
			WriteErrorJSON(w, err, 403)
			return
		}
		log.Error("failed to login user", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
