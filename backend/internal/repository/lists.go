package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/programmer8760/demonlist/backend/internal/domain"
)

func (r *postgresListsRepository) Create(ctx context.Context, list *domain.List) error {
	return r.db.QueryRowContext(
		ctx,
		"INSERT INTO lists (name) VALUES ($1) RETURNING id",
		list.Name,
	).Scan(&list.ID)
}

func (r *postgresListsRepository) GetAll(ctx context.Context) ([]domain.List, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM lists ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := []domain.List{}
	for rows.Next() {
		var list domain.List
		if err := rows.Scan(&list.ID, &list.Name); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, rows.Err()
}

func (r *postgresListsRepository) GetByID(ctx context.Context, id uint) (domain.List, error) {
	var list domain.List
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name FROM lists WHERE id = $1",
		id,
	).Scan(&list.ID, &list.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.List{}, domain.ErrListNotFound
		}
		return domain.List{}, err
	}
	return list, nil
}

func (r *postgresListsRepository) Update(ctx context.Context, list domain.List) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE lists SET name = $1 WHERE id = $2",
		list.Name,
		list.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrListNotFound
	}
	return nil
}

func (r *postgresListsRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM lists WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrListNotFound
	}
	return nil
}
