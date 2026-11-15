package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (r *postgresPermissionsRepository) Create(ctx context.Context, perms *domain.UserPermissions) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO user_permissions (user_id, list_id, permissions) VALUES ($1, $2, $3)",
		perms.UserID,
		perms.ListID,
		perms.Permissions,
	)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			return domain.ErrUserPermissionsAlreadyExist
		case "23503": // foreign_key_violation
			if pqErr.Constraint == "user_permissions_list_id_fkey" {
				return domain.ErrListNotFound
			}
			return domain.ErrUserNotFound
		}
	}
	return err
}

func (r *postgresPermissionsRepository) GetByListID(ctx context.Context, listID uint) ([]domain.UserPermissions, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT user_id, list_id, permissions FROM user_permissions WHERE list_id = $1 ORDER BY user_id",
		listID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []domain.UserPermissions{}
	for rows.Next() {
		var p domain.UserPermissions
		if err := rows.Scan(&p.UserID, &p.ListID, &p.Permissions); err != nil {
			return nil, err
		}
		perms = append(perms, p)
	}
	return perms, rows.Err()
}

func (r *postgresPermissionsRepository) Get(ctx context.Context, userID uint, listID uint) (domain.UserPermissions, error) {
	var p domain.UserPermissions
	err := r.db.QueryRowContext(
		ctx,
		"SELECT user_id, list_id, permissions FROM user_permissions WHERE user_id = $1 AND list_id = $2",
		userID,
		listID,
	).Scan(&p.UserID, &p.ListID, &p.Permissions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.UserPermissions{}, domain.ErrUserPermissionsNotFound
		}
		return domain.UserPermissions{}, err
	}
	return p, nil
}

func (r *postgresPermissionsRepository) Update(ctx context.Context, perms domain.UserPermissions) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE user_permissions SET permissions = $1 WHERE user_id = $2 AND list_id = $3",
		perms.Permissions,
		perms.UserID,
		perms.ListID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserPermissionsNotFound
	}
	return nil
}

func (r *postgresPermissionsRepository) Delete(ctx context.Context, userID uint, listID uint) error {
	res, err := r.db.ExecContext(
		ctx,
		"DELETE FROM user_permissions WHERE user_id = $1 AND list_id = $2",
		userID,
		listID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrUserPermissionsNotFound
	}
	return nil
}
