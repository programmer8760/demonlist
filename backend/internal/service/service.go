package service

import (
	"context"
	"log/slog"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/repository"
)

type LevelsService interface {
	CreateLevel(ctx context.Context, level domain.Level) (domain.Level, error)
	GetLevelsByListID(ctx context.Context, listID uint) ([]domain.Level, error)
	GetLevelByID(ctx context.Context, id uint) (domain.Level, error)
	UpdateLevel(ctx context.Context, level domain.Level) (domain.Level, error)
	UpdateLevelPlacements(ctx context.Context, listID uint, levelIDs []uint) ([]domain.Level, error)
	DeleteLevel(ctx context.Context, id uint) error
}

type levelsService struct {
	repo      repository.LevelsRepository
	listsRepo repository.ListsRepository
	usersRepo repository.UsersRepository
	permsRepo repository.PermissionsRepository
	logger    *slog.Logger
}

func NewLevelsService(repo repository.LevelsRepository, listsRepo repository.ListsRepository, usersRepo repository.UsersRepository, permsRepo repository.PermissionsRepository, log *slog.Logger) LevelsService {
	return &levelsService{
		repo:      repo,
		listsRepo: listsRepo,
		usersRepo: usersRepo,
		permsRepo: permsRepo,
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
	UpdateUser(ctx context.Context, input UpdateUserInput) (domain.User, error)
	DeleteUser(ctx context.Context, id uint) error
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

type RecordsService interface {
	SubmitRecord(ctx context.Context, record domain.Record) (domain.Record, error)
	GetRecordsByLevelID(ctx context.Context, levelID uint, status string) ([]domain.Record, error)
	GetRecordsByUserID(ctx context.Context, userID uint, status string) ([]domain.Record, error)
	UpdateRecordStatus(ctx context.Context, id uint, status string) (domain.Record, error)
	DeleteRecord(ctx context.Context, id uint) error
}

type recordsService struct {
	repo       repository.RecordsRepository
	levelsRepo repository.LevelsRepository
	usersRepo  repository.UsersRepository
	permsRepo  repository.PermissionsRepository
	logger     *slog.Logger
}

func NewRecordsService(repo repository.RecordsRepository, levelsRepo repository.LevelsRepository, usersRepo repository.UsersRepository, permsRepo repository.PermissionsRepository, log *slog.Logger) RecordsService {
	return &recordsService{
		repo:       repo,
		levelsRepo: levelsRepo,
		usersRepo:  usersRepo,
		permsRepo:  permsRepo,
		logger:     log,
	}
}

type PermissionsService interface {
	AddUserPermissions(ctx context.Context, perms domain.UserPermissions) (domain.UserPermissions, error)
	GetListPermissions(ctx context.Context, listID uint) ([]domain.UserPermissions, error)
	GetUserPermissions(ctx context.Context, userID uint, listID uint) (domain.UserPermissions, error)
	UpdateUserPermissions(ctx context.Context, perms domain.UserPermissions) (domain.UserPermissions, error)
	DeleteUserPermissions(ctx context.Context, userID uint, listID uint) error
}

type permissionsService struct {
	repo      repository.PermissionsRepository
	usersRepo repository.UsersRepository
	logger    *slog.Logger
}

func NewPermissionsService(repo repository.PermissionsRepository, usersRepo repository.UsersRepository, log *slog.Logger) PermissionsService {
	return &permissionsService{
		repo:      repo,
		usersRepo: usersRepo,
		logger:    log,
	}
}
