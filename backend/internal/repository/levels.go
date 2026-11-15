package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/programmer8760/demonlist/backend/internal/domain"
)

const levelColumns = "id, list_id, name, COALESCE(description, ''), showcase_link, placement"

type rowScanner interface {
	Scan(dest ...any) error
}

func scanLevel(row rowScanner, level *domain.Level) error {
	var listID sql.NullInt64
	err := row.Scan(&level.ID, &listID, &level.Name, &level.Description, &level.ShowcaseLink, &level.Placement)
	if err != nil {
		return err
	}
	if listID.Valid {
		id := uint(listID.Int64)
		level.ListID = &id
	}
	return nil
}

func (r *postgresLevelsRepository) Create(ctx context.Context, level *domain.Level) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM levels WHERE list_id = $1",
		level.ListID,
	).Scan(&count)
	if err != nil {
		return err
	}
	if level.Placement < 1 || level.Placement > count+1 {
		return domain.ErrInvalidPlacement
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE levels SET placement = placement + 1 WHERE list_id = $1 AND placement >= $2",
		level.ListID,
		level.Placement,
	)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(
		ctx,
		"INSERT INTO levels (list_id, name, description, showcase_link, placement) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		level.ListID,
		level.Name,
		level.Description,
		level.ShowcaseLink,
		level.Placement,
	).Scan(&level.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" { // foreign_key_violation
			return domain.ErrListNotFound
		}
		return err
	}

	return tx.Commit()
}

func (r *postgresLevelsRepository) GetByListID(ctx context.Context, listID uint) ([]domain.Level, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT "+levelColumns+" FROM levels WHERE list_id = $1 ORDER BY placement",
		listID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	levels := []domain.Level{}
	for rows.Next() {
		var level domain.Level
		if err := scanLevel(rows, &level); err != nil {
			return nil, err
		}
		levels = append(levels, level)
	}
	return levels, rows.Err()
}

func (r *postgresLevelsRepository) GetByID(ctx context.Context, id uint) (domain.Level, error) {
	var level domain.Level
	err := scanLevel(r.db.QueryRowContext(
		ctx,
		"SELECT "+levelColumns+" FROM levels WHERE id = $1",
		id,
	), &level)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Level{}, domain.ErrLevelNotFound
		}
		return domain.Level{}, err
	}
	return level, nil
}

func (r *postgresLevelsRepository) Update(ctx context.Context, level domain.Level) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE levels SET name = $1, description = $2, showcase_link = $3 WHERE id = $4",
		level.Name,
		level.Description,
		level.ShowcaseLink,
		level.ID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrLevelNotFound
	}
	return nil
}

func (r *postgresLevelsRepository) UpdatePlacements(ctx context.Context, listID uint, levelIDs []uint) error {
	ids := make([]int64, len(levelIDs))
	for i, id := range levelIDs {
		ids[i] = int64(id)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM levels WHERE list_id = $1",
		listID,
	).Scan(&count)
	if err != nil {
		return err
	}
	if count != len(ids) {
		return domain.ErrInvalidLevelOrder
	}

	res, err := tx.ExecContext(
		ctx,
		`UPDATE levels SET placement = v.placement
		FROM unnest($2::bigint[]) WITH ORDINALITY AS v(id, placement)
		WHERE levels.id = v.id AND levels.list_id = $1`,
		listID,
		pq.Array(ids),
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != int64(len(ids)) {
		return domain.ErrInvalidLevelOrder
	}

	return tx.Commit()
}

func (r *postgresLevelsRepository) Delete(ctx context.Context, id uint) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var listID sql.NullInt64
	var placement int
	err = tx.QueryRowContext(
		ctx,
		"DELETE FROM levels WHERE id = $1 RETURNING list_id, placement",
		id,
	).Scan(&listID, &placement)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrLevelNotFound
		}
		return err
	}

	if listID.Valid {
		_, err = tx.ExecContext(
			ctx,
			"UPDATE levels SET placement = placement - 1 WHERE list_id = $1 AND placement > $2",
			listID.Int64,
			placement,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
