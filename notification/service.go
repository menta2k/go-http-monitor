package notification

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

func (s *Service) ListByMonitor(ctx context.Context, monitorID int64) ([]domain.Notification, error) {
	return s.repo.FindByMonitorID(ctx, monitorID)
}

func (s *Service) Get(ctx context.Context, id int64) (domain.Notification, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, monitorID int64, typ domain.NotificationType, target string, enabled bool) (domain.Notification, error) {
	n, err := domain.NewNotification(monitorID, typ, target, enabled)
	if err != nil {
		return domain.Notification{}, err
	}
	return s.repo.Create(ctx, n)
}

func (s *Service) Update(ctx context.Context, id int64, typ domain.NotificationType, target string, enabled bool) (domain.Notification, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Notification{}, err
	}

	validated, err := domain.NewNotification(existing.MonitorID, typ, target, enabled)
	if err != nil {
		return domain.Notification{}, err
	}

	updated := domain.Notification{
		ID:        existing.ID,
		MonitorID: existing.MonitorID,
		Type:      validated.Type,
		Target:    validated.Target,
		Enabled:   validated.Enabled,
		CreatedAt: existing.CreatedAt,
	}

	return s.repo.Update(ctx, updated)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) FindEnabledByMonitor(ctx context.Context, monitorID int64) ([]domain.Notification, error) {
	return s.repo.FindEnabledByMonitorID(ctx, monitorID)
}
