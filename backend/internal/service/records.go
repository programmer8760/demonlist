package service

import (
	"context"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func validateVideoLink(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	link, err := url.ParseRequestURI(raw)
	if err != nil || (link.Scheme != "http" && link.Scheme != "https") || link.Host == "" || utf8.RuneCountInString(raw) > 255 {
		return "", domain.BadRequest{Err: domain.ErrInvalidVideoLink}
	}
	return raw, nil
}

func validateRecordStatus(status string) error {
	switch status {
	case domain.RecordStatusPending, domain.RecordStatusAccepted, domain.RecordStatusRejected:
		return nil
	}
	return domain.BadRequest{Err: domain.ErrInvalidRecordStatus}
}

func (s *recordsService) requireRecordPermissions(ctx context.Context, level domain.Level) error {
	if level.ListID == nil {
		return requireOwner(ctx, s.usersRepo)
	}
	return requirePermissions(ctx, s.usersRepo, s.permsRepo, *level.ListID, domain.PermChangeRecordStatus)
}

func (s *recordsService) SubmitRecord(ctx context.Context, record domain.Record) (domain.Record, error) {
	user, err := currentUser(ctx, s.usersRepo)
	if err != nil {
		return domain.Record{}, err
	}
	record.UserID = user.ID

	if record.Percentage < 1 || record.Percentage > 100 {
		return domain.Record{}, domain.BadRequest{Err: domain.ErrInvalidPercentage}
	}
	link, err := validateVideoLink(record.VideoLink)
	if err != nil {
		return domain.Record{}, err
	}
	record.VideoLink = link

	if err := s.repo.Create(ctx, &record); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *recordsService) GetRecordsByLevelID(ctx context.Context, levelID uint, status string) ([]domain.Record, error) {
	if status == "" {
		status = domain.RecordStatusAccepted
	}
	if err := validateRecordStatus(status); err != nil {
		return nil, err
	}

	level, err := s.levelsRepo.GetByID(ctx, levelID)
	if err != nil {
		return nil, err
	}
	if status != domain.RecordStatusAccepted {
		if err := s.requireRecordPermissions(ctx, level); err != nil {
			return nil, err
		}
	}
	return s.repo.GetByLevelID(ctx, levelID, status)
}

func (s *recordsService) GetRecordsByUserID(ctx context.Context, userID uint, status string) ([]domain.Record, error) {
	if status == "" {
		status = domain.RecordStatusAccepted
	}
	if err := validateRecordStatus(status); err != nil {
		return nil, err
	}

	if _, err := s.usersRepo.GetUserByID(ctx, userID); err != nil {
		return nil, err
	}
	if status != domain.RecordStatusAccepted {
		caller, err := currentUser(ctx, s.usersRepo)
		if err != nil {
			return nil, err
		}
		if caller.Role != "owner" && caller.ID != userID {
			return nil, domain.ErrForbidden
		}
	}
	return s.repo.GetByUserID(ctx, userID, status)
}

func (s *recordsService) UpdateRecordStatus(ctx context.Context, id uint, status string) (domain.Record, error) {
	if _, err := currentUser(ctx, s.usersRepo); err != nil {
		return domain.Record{}, err
	}

	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Record{}, err
	}
	level, err := s.levelsRepo.GetByID(ctx, record.LevelID)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.requireRecordPermissions(ctx, level); err != nil {
		return domain.Record{}, err
	}
	if err := validateRecordStatus(status); err != nil {
		return domain.Record{}, err
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return domain.Record{}, err
	}
	record.Status = status
	return record, nil
}

func (s *recordsService) DeleteRecord(ctx context.Context, id uint) error {
	caller, err := currentUser(ctx, s.usersRepo)
	if err != nil {
		return err
	}

	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if caller.ID != record.UserID {
		level, err := s.levelsRepo.GetByID(ctx, record.LevelID)
		if err != nil {
			return err
		}
		if err := s.requireRecordPermissions(ctx, level); err != nil {
			return err
		}
	}
	return s.repo.Delete(ctx, id)
}
