package service

import (
	"context"
	"errors"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/repository"
)

func currentUser(ctx context.Context, usersRepo repository.UsersRepository) (domain.User, error) {
	userID, ok := domain.UserIDFromContext(ctx)
	if !ok {
		return domain.User{}, domain.ErrUnauthorized
	}
	user, err := usersRepo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.User{}, domain.ErrUnauthorized
		}
		return domain.User{}, err
	}
	return user, nil
}

func requireOwner(ctx context.Context, usersRepo repository.UsersRepository) error {
	user, err := currentUser(ctx, usersRepo)
	if err != nil {
		return err
	}
	if user.Role != "owner" {
		return domain.ErrForbidden
	}
	return nil
}
