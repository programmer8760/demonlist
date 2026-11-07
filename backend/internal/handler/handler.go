package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type Handler struct {
	routes *http.ServeMux
	logger *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	h := &Handler{
		routes: http.NewServeMux(),
		logger: log,
	}

	// h.routes.HandleFunc("GET /subscriptions", h.GetAllSubscriptions)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	reqID := r.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = uuid.NewString()
	}

	ctx := context.WithValue(r.Context(), domain.RequestIDKey, reqID)
	r = r.WithContext(ctx)
	ww := NewWrapWriter(w)

	h.routes.ServeHTTP(ww, r)

	h.logger.Debug("request",
		"request_id", reqID,
		"method", r.Method,
		"path", r.URL.Path,
		"query", r.URL.RawQuery,
		"status", ww.Status(),
		"duration_ms", time.Since(start).Milliseconds(),
	)
}

func WriteErrorJSON(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	resp := domain.ErrorResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(resp)
}
