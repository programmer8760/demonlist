package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/programmer8760/demonlist/backend/internal/domain"
)

const recordColumns = "id, user_id, level_id, status, percentage, video_link"

func scanRecord(row rowScanner, record *domain.Record) error {
	return row.Scan(&record.ID, &record.UserID, &record.LevelID, &record.Status, &record.Percentage, &record.VideoLink)
}

func (r *postgresRecordsRepository) Create(ctx context.Context, record *domain.Record) error {
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO records (user_id, level_id, percentage, video_link) VALUES ($1, $2, $3, $4) RETURNING id, status",
		record.UserID,
		record.LevelID,
		record.Percentage,
		record.VideoLink,
	).Scan(&record.ID, &record.Status)

	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" { // foreign_key_violation
		if pqErr.Constraint == "records_level_id_fkey" {
			return domain.ErrLevelNotFound
		}
		return domain.ErrUserNotFound
	}
	return err
}

func (r *postgresRecordsRepository) GetByID(ctx context.Context, id uint) (domain.Record, error) {
	var record domain.Record
	err := scanRecord(r.db.QueryRowContext(
		ctx,
		"SELECT "+recordColumns+" FROM records WHERE id = $1",
		id,
	), &record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Record{}, domain.ErrRecordNotFound
		}
		return domain.Record{}, err
	}
	return record, nil
}

func (r *postgresRecordsRepository) getMany(ctx context.Context, query string, args ...any) ([]domain.Record, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []domain.Record{}
	for rows.Next() {
		var record domain.Record
		if err := scanRecord(rows, &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, rows.Err()
}

func (r *postgresRecordsRepository) GetByLevelID(ctx context.Context, levelID uint, status string) ([]domain.Record, error) {
	return r.getMany(
		ctx,
		"SELECT "+recordColumns+" FROM records WHERE level_id = $1 AND status = $2 ORDER BY percentage DESC, id",
		levelID,
		status,
	)
}

func (r *postgresRecordsRepository) GetByUserID(ctx context.Context, userID uint, status string) ([]domain.Record, error) {
	return r.getMany(
		ctx,
		"SELECT "+recordColumns+" FROM records WHERE user_id = $1 AND status = $2 ORDER BY id",
		userID,
		status,
	)
}

func (r *postgresRecordsRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE records SET status = $1 WHERE id = $2",
		status,
		id,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}

func (r *postgresRecordsRepository) Delete(ctx context.Context, id uint) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM records WHERE id = $1", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}
