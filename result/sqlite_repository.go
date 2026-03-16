package result

import (
	"context"
	"database/sql"
	"time"

	"github.com/sko/go-http-monitor/domain"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Create(ctx context.Context, cr domain.CheckResult) (domain.CheckResult, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO check_results (monitor_id, status_code, response_time_ms, body_matched, error, checked_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		cr.MonitorID, cr.StatusCode, cr.ResponseTimeMs, cr.BodyMatched, cr.Error,
		cr.CheckedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return domain.CheckResult{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.CheckResult{}, err
	}

	return domain.CheckResult{
		ID:             id,
		MonitorID:      cr.MonitorID,
		StatusCode:     cr.StatusCode,
		ResponseTimeMs: cr.ResponseTimeMs,
		BodyMatched:    cr.BodyMatched,
		Error:          cr.Error,
		CheckedAt:      cr.CheckedAt,
	}, nil
}

func (r *SQLiteRepository) FindByMonitorID(ctx context.Context, monitorID int64, limit, offset int) ([]domain.CheckResult, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, monitor_id, status_code, response_time_ms, body_matched, error, checked_at
		 FROM check_results WHERE monitor_id = ? ORDER BY checked_at DESC LIMIT ? OFFSET ?`,
		monitorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.CheckResult
	for rows.Next() {
		cr, err := scanResult(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, cr)
	}
	return results, rows.Err()
}

func (r *SQLiteRepository) CountByMonitorID(ctx context.Context, monitorID int64) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM check_results WHERE monitor_id = ?`, monitorID).Scan(&count)
	return count, err
}

func (r *SQLiteRepository) FindLatestByMonitorID(ctx context.Context, monitorID int64) (domain.CheckResult, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, monitor_id, status_code, response_time_ms, body_matched, error, checked_at
		 FROM check_results WHERE monitor_id = ? ORDER BY checked_at DESC LIMIT 1`,
		monitorID)

	var cr domain.CheckResult
	var checkedAt string
	var bodyMatched sql.NullBool
	err := row.Scan(&cr.ID, &cr.MonitorID, &cr.StatusCode, &cr.ResponseTimeMs, &bodyMatched, &cr.Error, &checkedAt)
	if err != nil {
		return domain.CheckResult{}, err
	}
	if bodyMatched.Valid {
		cr.BodyMatched = &bodyMatched.Bool
	}
	cr.CheckedAt = parseTime(checkedAt)
	return cr, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanResult(s scanner) (domain.CheckResult, error) {
	var cr domain.CheckResult
	var checkedAt string
	var bodyMatched sql.NullBool
	err := s.Scan(&cr.ID, &cr.MonitorID, &cr.StatusCode, &cr.ResponseTimeMs, &bodyMatched, &cr.Error, &checkedAt)
	if err != nil {
		return domain.CheckResult{}, err
	}
	if bodyMatched.Valid {
		cr.BodyMatched = &bodyMatched.Bool
	}
	cr.CheckedAt = parseTime(checkedAt)
	return cr, nil
}

func parseTime(s string) time.Time {
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05Z"} {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}
