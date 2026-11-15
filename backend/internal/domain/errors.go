package domain

import "errors"

type BadRequest struct {
	Err error
}

func (b BadRequest) Error() string {
	return b.Err.Error()
}

var (
	ErrUserAlreadyExists = errors.New("specified username is already taken by other user")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidID         = errors.New("invalid id")
	ErrInvalidListName   = errors.New("invalid list name")
	ErrListNotFound      = errors.New("list not found")
)

type ErrorResponse struct {
	Error string `json:"error"`
}
