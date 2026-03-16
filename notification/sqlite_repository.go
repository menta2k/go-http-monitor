package notification

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

func (r *SQLiteRepository) FindByMonitorID(ctx context.Context, monitorID int64) ([]domain.Notification, error) {
	return r.queryMultiple(ctx,
		`SELECT id, monitor_id, type, target, enabled, created_at, updated_at
		 FROM notifications WHERE monitor_id = ? ORDER BY id`, monitorID)
}

func (r *SQLiteRepository) FindEnabledByMonitorID(ctx context.Context, monitorID int64) ([]domain.Notification, error) {
	return r.queryMultiple(ctx,
		`SELECT id, monitor_id, type, target, enabled, created_at, updated_at
		 FROM notifications WHERE monitor_id = ? AND enabled = 1 ORDER BY id`, monitorID)
}

func (r *SQLiteRepository) FindByID(ctx context.Context, id int64) (domain.Notification, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, monitor_id, type, target, enabled, created_at, updated_at
		 FROM notifications WHERE id = ?`, id)
	return scanNotification(row)
}

func (r *SQLiteRepository) Create(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO notifications (monitor_id, type, target, enabled, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		n.MonitorID, string(n.Type), n.Target, n.Enabled,
		n.CreatedAt.UTC().Format(time.RFC3339), n.UpdatedAt.UTC().Format(time.RFC3339))
	if err != nil {
		return domain.Notification{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Notification{}, err
	}

	return domain.Notification{
		ID:        id,
		MonitorID: n.MonitorID,
		Type:      n.Type,
		Target:    n.Target,
		Enabled:   n.Enabled,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}, nil
}

func (r *SQLiteRepository) Update(ctx context.Context, n domain.Notification) (domain.Notification, error) {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET type = ?, target = ?, enabled = ?, updated_at = ? WHERE id = ?`,
		string(n.Type), n.Target, n.Enabled, now.Format(time.RFC3339), n.ID)
	if err != nil {
		return domain.Notification{}, err
	}

	return domain.Notification{
		ID:        n.ID,
		MonitorID: n.MonitorID,
		Type:      n.Type,
		Target:    n.Target,
		Enabled:   n.Enabled,
		CreatedAt: n.CreatedAt,
		UpdatedAt: now,
	}, nil
}

func (r *SQLiteRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM notifications WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

func (r *SQLiteRepository) queryMultiple(ctx context.Context, query string, args ...any) ([]domain.Notification, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []domain.Notification
	for rows.Next() {
		n, err := scanNotificationRow(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanNotification(s scanner) (domain.Notification, error) {
	n, err := scanNotificationRow(s)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Notification{}, domain.ErrNotificationNotFound
	}
	return n, err
}

func scanNotificationRow(s scanner) (domain.Notification, error) {
	var n domain.Notification
	var typ, createdAt, updatedAt string
	err := s.Scan(&n.ID, &n.MonitorID, &typ, &n.Target, &n.Enabled, &createdAt, &updatedAt)
	if err != nil {
		return domain.Notification{}, err
	}
	n.Type = domain.NotificationType(typ)
	n.CreatedAt = parseTime(createdAt)
	n.UpdatedAt = parseTime(updatedAt)
	return n, nil
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
