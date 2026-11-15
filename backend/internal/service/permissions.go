package service

import (
	"context"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func validatePermissions(perms uint16) error {
	if perms > domain.PermAll {
		return domain.BadRequest{Err: domain.ErrInvalidPermissions}
	}
	return nil
}

func (s *permissionsService) AddUserPermissions(ctx context.Context, perms domain.UserPermissions) (domain.UserPermissions, error) {
	if err := requireOwner(ctx, s.usersRepo); err != nil {
		return domain.UserPermissions{}, err
	}
	if err := validatePermissions(perms.Permissions); err != nil {
		return domain.UserPermissions{}, err
	}

	user, err := s.usersRepo.GetUserByID(ctx, perms.UserID)
	if err != nil {
		return domain.UserPermissions{}, err
	}
	if user.Role != "moderator" {
		return domain.UserPermissions{}, domain.BadRequest{Err: domain.ErrUserNotModerator}
	}

	if err := s.repo.Create(ctx, &perms); err != nil {
		return domain.UserPermissions{}, err
	}
	return perms, nil
}

func (s *permissionsService) GetListPermissions(ctx context.Context, listID uint) ([]domain.UserPermissions, error) {
	if err := requireOwner(ctx, s.usersRepo); err != nil {
		return nil, err
	}
	return s.repo.GetByListID(ctx, listID)
}

func (s *permissionsService) GetUserPermissions(ctx context.Context, userID uint, listID uint) (domain.UserPermissions, error) {
	caller, err := currentUser(ctx, s.usersRepo)
	if err != nil {
		return domain.UserPermissions{}, err
	}
	if caller.Role != "owner" && caller.ID != userID {
		return domain.UserPermissions{}, domain.ErrForbidden
	}

	target := caller
	if caller.ID != userID {
		target, err = s.usersRepo.GetUserByID(ctx, userID)
		if err != nil {
			return domain.UserPermissions{}, err
		}
	}
	if target.Role == "owner" {
		return domain.UserPermissions{UserID: userID, ListID: listID, Permissions: domain.PermAll}, nil
	}

	return s.repo.Get(ctx, userID, listID)
}

func (s *permissionsService) UpdateUserPermissions(ctx context.Context, perms domain.UserPermissions) (domain.UserPermissions, error) {
	if err := requireOwner(ctx, s.usersRepo); err != nil {
		return domain.UserPermissions{}, err
	}
	if err := validatePermissions(perms.Permissions); err != nil {
		return domain.UserPermissions{}, err
	}
	if err := s.repo.Update(ctx, perms); err != nil {
		return domain.UserPermissions{}, err
	}
	return perms, nil
}

func (s *permissionsService) DeleteUserPermissions(ctx context.Context, userID uint, listID uint) error {
	if err := requireOwner(ctx, s.usersRepo); err != nil {
		return err
	}
	return s.repo.Delete(ctx, userID, listID)
}
