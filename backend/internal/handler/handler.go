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
	routes      *http.ServeMux
	lists       service.ListsService
	levels      service.LevelsService
	users       service.UsersService
	records     service.RecordsService
	permissions service.PermissionsService
	logger      *slog.Logger
}

func NewHandler(log *slog.Logger, lists service.ListsService, levels service.LevelsService, users service.UsersService, records service.RecordsService, permissions service.PermissionsService) *Handler {
	h := &Handler{
		routes:      http.NewServeMux(),
		lists:       lists,
		levels:      levels,
		users:       users,
		records:     records,
		permissions: permissions,
		logger:      log,
	}

	h.routes.HandleFunc("POST /register", h.RegisterUser)
	h.routes.HandleFunc("POST /login", h.LoginUser)

	h.routes.HandleFunc("GET /lists", h.GetLists)
	h.routes.HandleFunc("GET /lists/{id}", h.GetListByID)
	h.routes.Handle("POST /lists", h.Auth(http.HandlerFunc(h.CreateList)))
	h.routes.Handle("PUT /lists/{id}", h.Auth(http.HandlerFunc(h.UpdateList)))
	h.routes.Handle("DELETE /lists/{id}", h.Auth(http.HandlerFunc(h.DeleteList)))

	h.routes.Handle("POST /lists/{id}/permissions", h.Auth(http.HandlerFunc(h.CreatePermissions)))
	h.routes.Handle("GET /lists/{id}/permissions", h.Auth(http.HandlerFunc(h.GetListPermissions)))
	h.routes.Handle("GET /lists/{id}/permissions/{user_id}", h.Auth(http.HandlerFunc(h.GetUserPermissions)))
	h.routes.Handle("PUT /lists/{id}/permissions/{user_id}", h.Auth(http.HandlerFunc(h.UpdatePermissions)))
	h.routes.Handle("DELETE /lists/{id}/permissions/{user_id}", h.Auth(http.HandlerFunc(h.DeletePermissions)))

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
