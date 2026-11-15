package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func validateLevel(level domain.Level) (domain.Level, error) {
	level.Name = strings.TrimSpace(level.Name)
	if level.Name == "" || utf8.RuneCountInString(level.Name) > 20 {
		return domain.Level{}, domain.BadRequest{Err: domain.ErrInvalidLevelName}
	}

	level.Description = strings.TrimSpace(level.Description)
	if utf8.RuneCountInString(level.Description) > 255 {
		return domain.Level{}, domain.BadRequest{Err: domain.ErrInvalidLevelDescription}
	}

	level.ShowcaseLink = strings.TrimSpace(level.ShowcaseLink)
	link, err := url.ParseRequestURI(level.ShowcaseLink)
	if err != nil || (link.Scheme != "http" && link.Scheme != "https") || link.Host == "" {
		return domain.Level{}, domain.BadRequest{Err: domain.ErrInvalidShowcaseLink}
	}

	return level, nil
}

func (s *levelsService) requireLevelPermissions(ctx context.Context, level domain.Level, required uint16) error {
	if level.ListID == nil {
		return requireOwner(ctx, s.usersRepo)
	}
	return requirePermissions(ctx, s.usersRepo, s.permsRepo, *level.ListID, required)
}

func (s *levelsService) CreateLevel(ctx context.Context, level domain.Level) (domain.Level, error) {
	if err := requirePermissions(ctx, s.usersRepo, s.permsRepo, *level.ListID, domain.PermAddLevels); err != nil {
		return domain.Level{}, err
	}

	level, err := validateLevel(level)
	if err != nil {
		return domain.Level{}, err
	}
	if level.Placement < 1 {
		return domain.Level{}, domain.BadRequest{Err: domain.ErrInvalidPlacement}
	}

	if err := s.repo.Create(ctx, &level); err != nil {
		if errors.Is(err, domain.ErrInvalidPlacement) {
			return domain.Level{}, domain.BadRequest{Err: err}
		}
		return domain.Level{}, err
	}
	return level, nil
}

func (s *levelsService) GetLevelsByListID(ctx context.Context, listID uint) ([]domain.Level, error) {
	if _, err := s.listsRepo.GetByID(ctx, listID); err != nil {
		return nil, err
	}
	return s.repo.GetByListID(ctx, listID)
}

func (s *levelsService) GetLevelByID(ctx context.Context, id uint) (domain.Level, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *levelsService) UpdateLevel(ctx context.Context, input domain.Level) (domain.Level, error) {
	level, err := s.repo.GetByID(ctx, input.ID)
	if err != nil {
		return domain.Level{}, err
	}
	if err := s.requireLevelPermissions(ctx, level, domain.PermEditLevels); err != nil {
		return domain.Level{}, err
	}

	input.ListID = level.ListID
	input.Placement = level.Placement
	level, err = validateLevel(input)
	if err != nil {
		return domain.Level{}, err
	}

	if err := s.repo.Update(ctx, level); err != nil {
		return domain.Level{}, err
	}
	return level, nil
}

func (s *levelsService) UpdateLevelPlacements(ctx context.Context, listID uint, levelIDs []uint) ([]domain.Level, error) {
	if err := requirePermissions(ctx, s.usersRepo, s.permsRepo, listID, domain.PermChangePlacements); err != nil {
		return nil, err
	}
	if _, err := s.listsRepo.GetByID(ctx, listID); err != nil {
		return nil, err
	}

	if err := s.repo.UpdatePlacements(ctx, listID, levelIDs); err != nil {
		if errors.Is(err, domain.ErrInvalidLevelOrder) {
			return nil, domain.BadRequest{Err: err}
		}
		return nil, err
	}
	return s.repo.GetByListID(ctx, listID)
}

func (s *levelsService) DeleteLevel(ctx context.Context, id uint) error {
	level, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.requireLevelPermissions(ctx, level, domain.PermDeleteLevels); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
