package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
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

func (r *postgresUsersRepository) Update(ctx context.Context, user domain.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(
		ctx,
		"UPDATE users SET username = $1, password_hash = $2, role = $3 WHERE id = $4",
		user.Username,
		user.Password,
		user.Role,
		user.ID,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" { // unique_violation
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}

	if user.Role != "moderator" {
		_, err = tx.ExecContext(ctx, "DELETE FROM user_permissions WHERE user_id = $1", user.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *postgresUsersRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}
