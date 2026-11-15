package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/utils"
)

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || tokenString == "" {
			WriteErrorJSON(w, domain.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		claims, err := utils.CheckJWT(tokenString)
		if err != nil {
			WriteErrorJSON(w, domain.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		userID, err := utils.GetUserIDFromJWTClaims(claims)
		if err != nil {
			WriteErrorJSON(w, domain.ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
