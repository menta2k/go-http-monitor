package database

import "database/sql"

func Migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS monitors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			url TEXT NOT NULL,
			expected_status INTEGER NOT NULL DEFAULT 200,
			body_contains TEXT NOT NULL DEFAULT '',
			interval_seconds INTEGER NOT NULL DEFAULT 60,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS check_results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			monitor_id INTEGER NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
			status_code INTEGER NOT NULL DEFAULT 0,
			response_time_ms INTEGER NOT NULL DEFAULT 0,
			body_matched BOOLEAN,
			error TEXT NOT NULL DEFAULT '',
			checked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_check_results_monitor_checked
			ON check_results(monitor_id, checked_at DESC)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			monitor_id INTEGER NOT NULL REFERENCES monitors(id) ON DELETE CASCADE,
			type TEXT NOT NULL CHECK(type IN ('email', 'slack')),
			target TEXT NOT NULL,
			enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_notifications_monitor
			ON notifications(monitor_id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
}
