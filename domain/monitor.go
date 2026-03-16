package domain

import (
	"fmt"
	"net/url"
	"time"
)

type Monitor struct {
	ID              int64     `json:"id"`
	URL             string    `json:"url"`
	ExpectedStatus  int       `json:"expected_status"`
	BodyContains    string    `json:"body_contains"`
	IntervalSeconds int       `json:"interval_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewMonitor(rawURL string, expectedStatus int, bodyContains string, intervalSeconds int) (Monitor, error) {
	if rawURL == "" {
		return Monitor{}, ErrInvalidURL
	}

	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return Monitor{}, fmt.Errorf("%w: must be a valid http/https URL", ErrInvalidURL)
	}

	if expectedStatus < 100 || expectedStatus > 599 {
		return Monitor{}, fmt.Errorf("%w: must be between 100 and 599", ErrInvalidStatusCode)
	}

	if intervalSeconds < 5 || intervalSeconds > 86400 {
		return Monitor{}, fmt.Errorf("%w: must be between 5 and 86400 seconds", ErrInvalidInterval)
	}

	now := time.Now().UTC()
	return Monitor{
		URL:             rawURL,
		ExpectedStatus:  expectedStatus,
		BodyContains:    bodyContains,
		IntervalSeconds: intervalSeconds,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
