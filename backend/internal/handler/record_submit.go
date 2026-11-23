package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type SubmitRecordRequest struct {
	Percentage int    `json:"percentage"`
	VideoLink  string `json:"video_link"`
}

func (h *Handler) SubmitRecord(w http.ResponseWriter, r *http.Request) {
	levelID, err := parseIDFromPath(r)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	var req SubmitRecordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		WriteErrorJSON(w, err, 400)
		return
	}

	ctx := r.Context()
	log := h.logger.With("request_id", ctx.Value(domain.RequestIDKey))
	log.Debug("request body", "input", req)

	resp, err := h.records.SubmitRecord(ctx, domain.Record{
		LevelID:    levelID,
		Percentage: req.Percentage,
		VideoLink:  req.VideoLink,
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
		if errors.Is(err, domain.ErrLevelNotFound) {
			WriteErrorJSON(w, err, 404)
			return
		}
		log.Error("failed to submit record", "err", err)
		WriteErrorJSON(w, errors.New("internal server error"), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
