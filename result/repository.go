package result

import (
	"context"

	"github.com/sko/go-http-monitor/domain"
)

type Repository interface {
	Create(ctx context.Context, r domain.CheckResult) (domain.CheckResult, error)
	FindByMonitorID(ctx context.Context, monitorID int64, limit, offset int) ([]domain.CheckResult, error)
	FindLatestByMonitorID(ctx context.Context, monitorID int64) (domain.CheckResult, error)
}
