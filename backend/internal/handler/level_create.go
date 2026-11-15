package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type CreateLevelRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	ShowcaseLink string `json:"showcase_link"`
	Placement    int    `json:"placement"`
}

func (h *Handler) CreateLevel(w http.ResponseWriter, r *http.Request) {
	listID, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	var req CreateLevelRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.levels.CreateLevel(ctx, domain.Level{
		ListID:       &listID,
		Name:         req.Name,
		Description:  req.Description,
		ShowcaseLink: req.ShowcaseLink,
		Placement:    req.Placement,
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
		if errors.Is(err, domain.ErrListNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to create level", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
