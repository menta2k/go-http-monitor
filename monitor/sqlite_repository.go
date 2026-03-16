package monitor

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sko/go-http-monitor/domain"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) FindAll(ctx context.Context) ([]domain.Monitor, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, url, expected_status, body_contains, interval_seconds, user_agent, created_at, updated_at
		 FROM monitors ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []domain.Monitor
	for rows.Next() {
		m, err := scanMonitor(rows)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, m)
	}
	return monitors, rows.Err()
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id int64) (domain.Monitor, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, url, expected_status, body_contains, interval_seconds, user_agent, created_at, updated_at
		 FROM monitors WHERE id = ?`, id)

	var m domain.Monitor
	var createdAt, updatedAt string
	err := row.Scan(&m.ID, &m.URL, &m.ExpectedStatus, &m.BodyContains, &m.IntervalSeconds, &m.UserAgent, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Monitor{}, domain.ErrMonitorNotFound
	}
	if err != nil {
		return domain.Monitor{}, err
	}
	m.CreatedAt = parseTime(createdAt)
	m.UpdatedAt = parseTime(updatedAt)
	return m, nil
}

func (r *SQLiteRepository) Create(ctx context.Context, m domain.Monitor) (domain.Monitor, error) {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO monitors (url, expected_status, body_contains, interval_seconds, user_agent, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		m.URL, m.ExpectedStatus, m.BodyContains, m.IntervalSeconds, m.UserAgent,
		m.CreatedAt.UTC().Format(time.RFC3339), m.UpdatedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return domain.Monitor{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return domain.Monitor{}, err
	}

	return domain.Monitor{
		ID:              id,
		URL:             m.URL,
		ExpectedStatus:  m.ExpectedStatus,
		BodyContains:    m.BodyContains,
		IntervalSeconds: m.IntervalSeconds,
		UserAgent:       m.UserAgent,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}, nil
}

func (r *SQLiteRepository) Update(ctx context.Context, m domain.Monitor) (domain.Monitor, error) {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx,
		`UPDATE monitors SET url = ?, expected_status = ?, body_contains = ?, interval_seconds = ?, user_agent = ?, updated_at = ?
		 WHERE id = ?`,
		m.URL, m.ExpectedStatus, m.BodyContains, m.IntervalSeconds, m.UserAgent,
		now.Format(time.RFC3339), m.ID)
	if err != nil {
		return domain.Monitor{}, err
	}

	return domain.Monitor{
		ID:              m.ID,
		URL:             m.URL,
		ExpectedStatus:  m.ExpectedStatus,
		BodyContains:    m.BodyContains,
		IntervalSeconds: m.IntervalSeconds,
		UserAgent:       m.UserAgent,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       now,
	}, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM monitors WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrMonitorNotFound
	}
	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanMonitor(s scanner) (domain.Monitor, error) {
	var m domain.Monitor
	var createdAt, updatedAt string
	err := s.Scan(&m.ID, &m.URL, &m.ExpectedStatus, &m.BodyContains, &m.IntervalSeconds, &m.UserAgent, &createdAt, &updatedAt)
	if err != nil {
		return domain.Monitor{}, err
	}
	m.CreatedAt = parseTime(createdAt)
	m.UpdatedAt = parseTime(updatedAt)
	return m, nil
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
