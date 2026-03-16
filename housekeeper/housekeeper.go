package housekeeper

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Housekeeper struct {
	db            *sql.DB
	interval      time.Duration
	retentionDays int
}

func New(db *sql.DB, interval time.Duration, retentionDays int) *Housekeeper {
	return &Housekeeper{
		db:            db,
		interval:      interval,
		retentionDays: retentionDays,
	}
}

func (h *Housekeeper) Run(ctx context.Context) {
	log.Printf("[housekeeper] started (interval=%s, retention=%dd)", h.interval, h.retentionDays)

	// Run immediately on startup
	h.cleanup(ctx)

	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[housekeeper] stopped")
			return
		case <-ticker.C:
			h.cleanup(ctx)
		}
	}
}

func (h *Housekeeper) cleanup(ctx context.Context) {
	cutoff := time.Now().UTC().AddDate(0, 0, -h.retentionDays).Format(time.RFC3339)

	result, err := h.db.ExecContext(ctx,
		`DELETE FROM check_results WHERE checked_at < ?`, cutoff)
	if err != nil {
		log.Printf("[housekeeper] failed to delete old results: %v", err)
		return
	}

	deleted, _ := result.RowsAffected()
	if deleted > 0 {
		log.Printf("[housekeeper] purged %d check results older than %d days", deleted, h.retentionDays)
	}
}
