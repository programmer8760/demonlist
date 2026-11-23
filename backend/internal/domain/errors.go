package domain

import "errors"

type BadRequest struct {
	Err error
}

func (b BadRequest) Error() string {
	return b.Err.Error()
}

var (
	ErrUserAlreadyExists           = errors.New("specified username is already taken by other user")
	ErrInvalidUsername             = errors.New("invalid username")
	ErrInvalidPassword             = errors.New("invalid password")
	ErrUserNotFound                = errors.New("user not found")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrForbidden                   = errors.New("forbidden")
	ErrInvalidID                   = errors.New("invalid id")
	ErrInvalidListName             = errors.New("invalid list name")
	ErrListNotFound                = errors.New("list not found")
	ErrInvalidPermissions          = errors.New("invalid permissions")
	ErrUserPermissionsNotFound     = errors.New("user permissions not found")
	ErrUserPermissionsAlreadyExist = errors.New("user already has permissions in this list")
	ErrUserNotModerator            = errors.New("user is not a moderator")
	ErrInvalidLevelName            = errors.New("invalid level name")
	ErrInvalidLevelDescription     = errors.New("invalid level description")
	ErrInvalidShowcaseLink         = errors.New("invalid showcase link")
	ErrInvalidPlacement            = errors.New("invalid placement")
	ErrInvalidLevelOrder           = errors.New("level order must contain every level of the list exactly once")
	ErrLevelNotFound               = errors.New("level not found")
	ErrInvalidPercentage           = errors.New("invalid percentage")
	ErrInvalidVideoLink            = errors.New("invalid video link")
	ErrInvalidRecordStatus         = errors.New("invalid record status")
	ErrRecordNotFound              = errors.New("record not found")
	ErrInvalidRole                 = errors.New("invalid role")
)

type ErrorResponse struct {
	Error string `json:"error"`
}
