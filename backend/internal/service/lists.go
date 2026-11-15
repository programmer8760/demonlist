package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (s *listsService) requireOwner(ctx context.Context) error {
	userID, ok := domain.UserIDFromContext(ctx)
	if !ok {
		return domain.ErrUnauthorized
	}
	user, err := s.usersRepo.GetUserByID(ctx, userID)
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

func validateListName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || utf8.RuneCountInString(name) > 32 {
		return "", domain.BadRequest{Err: domain.ErrInvalidListName}
	}
	return name, nil
}

func (s *listsService) CreateList(ctx context.Context, list domain.List) (domain.List, error) {
	if err := s.requireOwner(ctx); err != nil {
		return domain.List{}, err
	}

	name, err := validateListName(list.Name)
	if err != nil {
		return domain.List{}, err
	}
	list.Name = name

	if err := s.repo.Create(ctx, &list); err != nil {
		return domain.List{}, err
	}
	return list, nil
}

func (s *listsService) GetLists(ctx context.Context) ([]domain.List, error) {
	return s.repo.GetAll(ctx)
}

func (s *listsService) GetListByID(ctx context.Context, id uint) (domain.List, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *listsService) UpdateList(ctx context.Context, list domain.List) (domain.List, error) {
	if err := s.requireOwner(ctx); err != nil {
		return domain.List{}, err
	}

	name, err := validateListName(list.Name)
	if err != nil {
		return domain.List{}, err
	}
	list.Name = name

	if err := s.repo.Update(ctx, list); err != nil {
		return domain.List{}, err
	}
	return list, nil
}

func (s *listsService) DeleteList(ctx context.Context, id uint) error {
	if err := s.requireOwner(ctx); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
