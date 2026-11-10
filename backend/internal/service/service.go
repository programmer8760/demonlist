package service

import (
	"log/slog"

	"github.com/programmer8760/demonlist/backend/internal/repository"
)

type LevelsService interface{}

type levelsService struct {
	repo      repository.LevelsRepository
	listsRepo repository.ListsRepository
	logger    *slog.Logger
}

func NewLevelsService(repo repository.LevelsRepository, listsRepo repository.ListsRepository, log *slog.Logger) LevelsService {
	return &levelsService{
		repo:      repo,
		listsRepo: listsRepo,
		logger:    log,
	}
}

type ListsService interface{}

type listsService struct {
	repo   repository.ListsRepository
	logger *slog.Logger
}

func NewListsService(repo repository.ListsRepository, log *slog.Logger) ListsService {
	return &listsService{
		repo:   repo,
		logger: log,
	}
}

type UsersService interface{}

type usersService struct {
	repo        repository.UsersRepository
	recordsRepo repository.RecordsRepository
	logger      *slog.Logger
}

func NewUsersService(repo repository.UsersRepository, recordsRepo repository.RecordsRepository, log *slog.Logger) UsersService {
	return &usersService{
		repo:        repo,
		recordsRepo: recordsRepo,
		logger:      log,
	}
}

type RecordsService interface{}

type recordsService struct {
	repo      repository.RecordsRepository
	usersRepo repository.UsersRepository
	logger    *slog.Logger
}

func NewRecordsService(repo repository.RecordsRepository, usersRepo repository.UsersRepository, log *slog.Logger) RecordsService {
	return &recordsService{
		repo:      repo,
		usersRepo: usersRepo,
		logger:    log,
	}
}
