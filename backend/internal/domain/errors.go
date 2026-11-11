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
)

type ErrorResponse struct {
	Error string `json:"error"`
}
