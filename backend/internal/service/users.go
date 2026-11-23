package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/utils"
)

func validateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	if username == "" || utf8.RuneCountInString(username) > 32 {
		return "", domain.BadRequest{Err: domain.ErrInvalidUsername}
	}
	return username, nil
}

func validatePassword(password string) (string, error) {
	password = strings.TrimSpace(password)
	if password == "" || len(password) > 72 {
		return "", domain.BadRequest{Err: domain.ErrInvalidPassword}
	}
	return password, nil
}

type UpdateUserInput struct {
	ID       uint
	Username *string
	Password *string
	Role     *string
}

func (s *usersService) UpdateUser(ctx context.Context, input UpdateUserInput) (domain.User, error) {
	caller, err := currentUser(ctx, s.repo)
	if err != nil {
		return domain.User{}, err
	}
	if caller.Role != "owner" && caller.ID != input.ID {
		return domain.User{}, domain.ErrForbidden
	}

	user := caller
	if caller.ID != input.ID {
		user, err = s.repo.GetUserByID(ctx, input.ID)
		if err != nil {
			return domain.User{}, err
		}
	}

	if input.Username != nil {
		username, err := validateUsername(*input.Username)
		if err != nil {
			return domain.User{}, err
		}
		user.Username = username
	}

	if input.Password != nil {
		if caller.ID != user.ID {
			return domain.User{}, domain.ErrForbidden
		}
		password, err := validatePassword(*input.Password)
		if err != nil {
			return domain.User{}, err
		}
		hash, err := utils.HashPassword(password)
		if err != nil {
			return domain.User{}, err
		}
		user.Password = hash
	}

	if input.Role != nil {
		if caller.Role != "owner" || caller.ID == user.ID || user.Role == "owner" {
			return domain.User{}, domain.ErrForbidden
		}
		if *input.Role != "player" && *input.Role != "moderator" {
			return domain.User{}, domain.BadRequest{Err: domain.ErrInvalidRole}
		}
		user.Role = *input.Role
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return domain.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (s *usersService) DeleteUser(ctx context.Context, id uint) error {
	caller, err := currentUser(ctx, s.repo)
	if err != nil {
		return err
	}
	if caller.Role != "owner" && caller.ID != id {
		return domain.ErrForbidden
	}

	target := caller
	if caller.ID != id {
		target, err = s.repo.GetUserByID(ctx, id)
		if err != nil {
			return err
		}
	}
	if target.Role == "owner" {
		return domain.ErrForbidden
	}

	return s.repo.Delete(ctx, id)
}
