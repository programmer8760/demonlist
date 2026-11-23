package handler

import (
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (h *Handler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))

	err = h.records.DeleteRecord(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			WriteErrorJSON(w, err, 401)
			return
		}
		if errors.Is(err, domain.ErrForbidden) {
			WriteErrorJSON(w, err, 403)
			return
		}
		if errors.Is(err, domain.ErrRecordNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to delete record", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
