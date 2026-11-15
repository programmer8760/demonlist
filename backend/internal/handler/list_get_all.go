package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (h *Handler) GetLists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))

	resp, err := h.lists.GetLists(ctx)
	if err != nil {
		log.Error("failed to get lists", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
