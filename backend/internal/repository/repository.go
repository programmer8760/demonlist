package repository

import (
	"context"
	"database/sql"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

type LevelsRepository interface {
	Create(ctx context.Context, level *domain.Level) error
	GetByListID(ctx context.Context, listID uint) ([]domain.Level, error)
	GetByID(ctx context.Context, id uint) (domain.Level, error)
	Update(ctx context.Context, level domain.Level) error
	UpdatePlacements(ctx context.Context, listID uint, levelIDs []uint) error
	Delete(ctx context.Context, id uint) error
}

type postgresLevelsRepository struct {
	db *sql.DB
}

func NewLevelsRepository(db *sql.DB) LevelsRepository {
	return &postgresLevelsRepository{db: db}
}

type ListsRepository interface {
	Create(ctx context.Context, list *domain.List) error
	GetAll(ctx context.Context) ([]domain.List, error)
	GetByID(ctx context.Context, id uint) (domain.List, error)
	Update(ctx context.Context, list domain.List) error
	Delete(ctx context.Context, id uint) error
}

type postgresListsRepository struct {
	db *sql.DB
}

func NewListsRepository(db *sql.DB) ListsRepository {
	return &postgresListsRepository{db: db}
}

type UsersRepository interface {
	Register(ctx context.Context, user *domain.User) error
	GetUserByUsername(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uint) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id uint) error
}

type postgresUsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &postgresUsersRepository{db: db}
}

type RecordsRepository interface {
	Create(ctx context.Context, record *domain.Record) error
	GetByID(ctx context.Context, id uint) (domain.Record, error)
	GetByLevelID(ctx context.Context, levelID uint, status string) ([]domain.Record, error)
	GetByUserID(ctx context.Context, userID uint, status string) ([]domain.Record, error)
	UpdateStatus(ctx context.Context, id uint, status string) error
	Delete(ctx context.Context, id uint) error
}

type postgresRecordsRepository struct {
	db *sql.DB
}

func NewRecordsRepository(db *sql.DB) RecordsRepository {
	return &postgresRecordsRepository{db: db}
}

type PermissionsRepository interface {
	Create(ctx context.Context, perms *domain.UserPermissions) error
	GetByListID(ctx context.Context, listID uint) ([]domain.UserPermissions, error)
	Get(ctx context.Context, userID uint, listID uint) (domain.UserPermissions, error)
	Update(ctx context.Context, perms domain.UserPermissions) error
	Delete(ctx context.Context, userID uint, listID uint) error
}

type postgresPermissionsRepository struct {
	db *sql.DB
}

func NewPermissionsRepository(db *sql.DB) PermissionsRepository {
	return &postgresPermissionsRepository{db: db}
}
