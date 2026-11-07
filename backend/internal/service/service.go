package service

import (
	"log/slog"
)

type LevelsService interface{}

type levelsService struct {
	logger *slog.Logger
}

func NewLevelsService(log *slog.Logger) LevelsService {
	return &levelsService{
		logger: log,
	}
}

type ListsService interface{}

type listsService struct {
	logger *slog.Logger
}

func NewListsService(log *slog.Logger) ListsService {
	return &listsService{
		logger: log,
	}
}

type UsersService interface{}

type usersService struct {
	logger *slog.Logger
}

func NewUsersService(log *slog.Logger) UsersService {
	return &usersService{
		logger: log,
	}
}

type RecordsService interface{}

type recordsService struct {
	logger *slog.Logger
}

func NewRecordsService(log *slog.Logger) RecordsService {
	return &recordsService{
		logger: log,
	}
}
