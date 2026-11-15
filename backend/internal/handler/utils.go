package handler

import (
	"net/http"
	"strconv"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func parseIDFromPath(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, domain.BadRequest{Err: domain.ErrInvalidID}
	}
	return uint(id), nil
}

func parseUserIDFromPath(r *http.Request) (uint, error) {
	id, err := strconv.ParseUint(r.PathValue("user_id"), 10, 64)
	if err != nil || id == 0 {
		return 0, domain.BadRequest{Err: domain.ErrInvalidID}
	}
	return uint(id), nil
}
