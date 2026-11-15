package domain

import "context"

const (
	RequestIDKey string = "request_id"
	UserIDKey    string = "user_id"
)

func UserIDFromContext(ctx context.Context) (uint, bool) {
	id, ok := ctx.Value(UserIDKey).(uint)
	return id, ok
}
