package notification

import (
	"context"

	"github.com/sko/go-http-monitor/domain"
)

type Repository interface {
	FindByMonitorID(ctx context.Context, monitorID int64) ([]domain.Notification, error)
	FindByID(ctx context.Context, id int64) (domain.Notification, error)
	Create(ctx context.Context, n domain.Notification) (domain.Notification, error)
	Update(ctx context.Context, n domain.Notification) (domain.Notification, error)
	Delete(ctx context.Context, id int64) error
	FindEnabledByMonitorID(ctx context.Context, monitorID int64) ([]domain.Notification, error)
}
