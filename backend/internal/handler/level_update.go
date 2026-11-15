package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type UpdateLevelRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ShowcaseLink string `json:"showcase_link"`
}

func (h *Handler) UpdateLevel(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	var req UpdateLevelRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.levels.UpdateLevel(ctx, domain.Level{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		ShowcaseLink: req.ShowcaseLink,
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
		if errors.Is(err, domain.ErrLevelNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to update level", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
