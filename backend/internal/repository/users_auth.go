package repository

import (
	"context"
	"errors"

	"github.com/lib/pq"
	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (r *postgresUsersRepository) Register(ctx context.Context, user *domain.User) error {
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id, role",
		user.Username,
		user.Password,
	).Scan(&user.ID, &user.Role)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" { // unique_violation
			return domain.ErrUserAlreadyExists
		}
	}
	return err
}
