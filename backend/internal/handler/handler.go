package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/service"
)

type Handler struct {
	routes  *http.ServeMux
	lists   service.ListsService
	levels  service.LevelsService
	users   service.UsersService
	records service.RecordsService
	logger  *slog.Logger
}

func NewHandler(log *slog.Logger, lists service.ListsService, levels service.LevelsService, users service.UsersService, records service.RecordsService) *Handler {
	h := &Handler{
		routes:  http.NewServeMux(),
		lists:   lists,
		levels:  levels,
		users:   users,
		records: records,
		logger:  log,
	}

	h.routes.HandleFunc("POST /register", h.RegisterUser)
	h.routes.HandleFunc("POST /login", h.LoginUser)

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
