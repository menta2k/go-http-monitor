package domain

import "time"

type CheckResult struct {
	ID             int64     `json:"id"`
	MonitorID      int64     `json:"monitor_id"`
	StatusCode     int       `json:"status_code"`
	ResponseTimeMs int64     `json:"response_time_ms"`
	BodyMatched    *bool     `json:"body_matched"`
	Error          string    `json:"error,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}
