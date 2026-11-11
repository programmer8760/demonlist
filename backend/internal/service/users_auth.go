package service

import (
	"context"
	"strings"

	"github.com/programmer8760/demonlist/backend/internal/domain"
	"github.com/programmer8760/demonlist/backend/internal/utils"
)

type RegisterResponse struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

func (s *usersService) RegisterUser(ctx context.Context, user domain.User) (RegisterResponse, error) {
	if user.Username = strings.TrimSpace(user.Username); user.Username == "" {
		return RegisterResponse{}, domain.BadRequest{Err: domain.ErrInvalidUsername}
	}
	if user.Password = strings.TrimSpace(user.Password); user.Password == "" {
		return RegisterResponse{}, domain.BadRequest{Err: domain.ErrInvalidPassword}
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return RegisterResponse{}, err
	}
	user.Password = hash

	if err := s.repo.Register(ctx, &user); err != nil {
		return RegisterResponse{}, err
	}

	user.Password = ""

	token, err := utils.IssueJWT(user.ID)
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{user, token}, nil
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

func (s *usersService) LoginUser(ctx context.Context, attempt domain.User) (LoginUserResponse, error) {
	if attempt.Username = strings.TrimSpace(attempt.Username); attempt.Username == "" {
		return LoginUserResponse{}, domain.BadRequest{Err: domain.ErrInvalidUsername}
	}

	user, err := s.repo.GetUserByUsername(ctx, attempt)
	if err != nil {
		return LoginUserResponse{}, err
	}

	err = utils.CheckPassword(user.Password, user.Password)
	if err != nil {
		return LoginUserResponse{}, domain.ErrInvalidPassword
	}

	token, err := utils.IssueJWT(user.ID)
	if err != nil {
		return LoginUserResponse{}, err
	}

	return LoginUserResponse{token}, nil
}
