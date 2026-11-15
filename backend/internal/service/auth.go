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

func requirePermissions(ctx context.Context, usersRepo repository.UsersRepository, permsRepo repository.PermissionsRepository, listID uint, required uint16) error {
	user, err := currentUser(ctx, usersRepo)
	if err != nil {
		return err
	}
	if user.Role == "owner" {
		return nil
	}
	if user.Role != "moderator" {
		return domain.ErrForbidden
	}

	perms, err := permsRepo.Get(ctx, user.ID, listID)
	if err != nil {
		if errors.Is(err, domain.ErrUserPermissionsNotFound) {
			return domain.ErrForbidden
		}
		return err
	}
	if !domain.HasPermissions(perms.Permissions, required) {
		return domain.ErrForbidden
	}
	return nil
}
