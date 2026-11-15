package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (r *postgresUsersRepository) GetUserByID(ctx context.Context, id uint) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRowContext(
		ctx,
		"SELECT * FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}
