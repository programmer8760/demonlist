package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (h *Handler) GetLevelsByListID(w http.ResponseWriter, r *http.Request) {
	listID, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))

	resp, err := h.levels.GetLevelsByListID(ctx, listID)
	if err != nil {
		if errors.Is(err, domain.ErrListNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to get levels", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
