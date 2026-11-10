package repository

import (
	"database/sql"
)

type LevelsRepository interface{}

type postgresLevelsRepository struct {
	db *sql.DB
}

func NewLevelsRepository(db *sql.DB) LevelsRepository {
	return &postgresLevelsRepository{db: db}
}

type ListsRepository interface{}

type postgresListsRepository struct {
	db *sql.DB
}

func NewListsRepository(db *sql.DB) ListsRepository {
	return &postgresListsRepository{db: db}
}

type UsersRepository interface{}

type postgresUsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &postgresUsersRepository{db: db}
}

type RecordsRepository interface{}

type postgresRecordsRepository struct {
	db *sql.DB
}

func NewRecordsRepository(db *sql.DB) RecordsRepository {
	return &postgresRecordsRepository{db: db}
}

type PermissionsRepository interface{}

type postgresPermissionsRepository struct {
	db *sql.DB
}

func NewPermissionsRepository(db *sql.DB) PermissionsRepository {
	return &postgresPermissionsRepository{db: db}
}
