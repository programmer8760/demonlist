package service

import (
	"context"
	"log/slog"

	"github.com/programmer8760/demonlist/backend/internal/domain"
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

type ListsService interface {
	CreateList(ctx context.Context, list domain.List) (domain.List, error)
	GetLists(ctx context.Context) ([]domain.List, error)
	GetListByID(ctx context.Context, id uint) (domain.List, error)
	UpdateList(ctx context.Context, list domain.List) (domain.List, error)
	DeleteList(ctx context.Context, id uint) error
}

type listsService struct {
	repo      repository.ListsRepository
	usersRepo repository.UsersRepository
	logger    *slog.Logger
}

func NewListsService(repo repository.ListsRepository, usersRepo repository.UsersRepository, log *slog.Logger) ListsService {
	return &listsService{
		repo:      repo,
		usersRepo: usersRepo,
		logger:    log,
	}
}

type UsersService interface {
	RegisterUser(ctx context.Context, user domain.User) (RegisterResponse, error)
	LoginUser(ctx context.Context, user domain.User) (LoginUserResponse, error)
}

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
