package domain

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"time"
)

type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationSlack NotificationType = "slack"
)

type Notification struct {
	ID        int64            `json:"id"`
	MonitorID int64            `json:"monitor_id"`
	Type      NotificationType `json:"type"`
	Target    string           `json:"target"`
	Enabled   bool             `json:"enabled"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotificationType = errors.New("invalid notification type: must be 'email' or 'slack'")
	ErrInvalidNotificationTarget = errors.New("invalid notification target")
)

func NewNotification(monitorID int64, typ NotificationType, target string, enabled bool) (Notification, error) {
	switch typ {
	case NotificationEmail:
		if _, err := mail.ParseAddress(target); err != nil {
			return Notification{}, fmt.Errorf("%w: invalid email address", ErrInvalidNotificationTarget)
		}
	case NotificationSlack:
		u, err := url.Parse(target)
		if err != nil || u.Scheme != "https" || u.Host == "" {
			return Notification{}, fmt.Errorf("%w: must be a valid https webhook URL", ErrInvalidNotificationTarget)
		}
	default:
		return Notification{}, ErrInvalidNotificationType
	}

	if monitorID <= 0 {
		return Notification{}, ErrMonitorNotFound
	}

	now := time.Now().UTC()
	return Notification{
		MonitorID: monitorID,
		Type:      typ,
		Target:    target,
		Enabled:   enabled,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
