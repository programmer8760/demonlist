package domain

import "errors"

type BadRequest struct {
	Err error
}

func (b BadRequest) Error() string {
	return b.Err.Error()
}

var ErrTestError = errors.New("test error")

type ErrorResponse struct {
	Error string `json:"error"`
}
