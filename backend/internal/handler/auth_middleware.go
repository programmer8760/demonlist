package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/utils"
)

func userIDFromRequest(r *http.Request) (uint, error) {
	tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok || tokenString == "" {
		return 0, domain.ErrUnauthorized
	}

	claims, err := utils.CheckJWT(tokenString)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}

	userID, err := utils.GetUserIDFromJWTClaims(claims)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	return userID, nil
}

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := userIDFromRequest(r)
		if err != nil {
			WriteErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			next.ServeHTTP(w, r)
			return
		}
		h.Auth(next).ServeHTTP(w, r)
	})
}
