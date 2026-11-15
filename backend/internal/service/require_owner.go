package service

import (
	"context"
	"errors"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/repository"
)

func requireOwner(ctx context.Context, usersRepo repository.UsersRepository) error {
	userID, ok := domain.UserIDFromContext(ctx)
	if !ok {
		return domain.ErrUnauthorized
	}
	user, err := usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUnauthorized
		}
		return err
	}
	if user.Role != "owner" {
		return domain.ErrForbidden
	}
	return nil
}
