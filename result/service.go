package result

import (
	"context"

	"github.com/sko/go-http-monitor/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Latest(ctx context.Context, monitorID int64) (domain.CheckResult, error) {
	return s.repo.FindLatestByMonitorID(ctx, monitorID)
}

func (s *Service) History(ctx context.Context, monitorID int64, limit, offset int) ([]domain.CheckResult, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.FindByMonitorID(ctx, monitorID, limit, offset)
}
