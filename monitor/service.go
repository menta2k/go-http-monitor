package monitor

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

func (s *Service) List(ctx context.Context) ([]domain.Monitor, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) Get(ctx context.Context, id int64) (domain.Monitor, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, url string, expectedStatus int, bodyContains string, intervalSeconds int, userAgent string) (domain.Monitor, error) {
	m, err := domain.NewMonitor(url, expectedStatus, bodyContains, intervalSeconds, userAgent)
	if err != nil {
		return domain.Monitor{}, err
	}
	return s.repo.Create(ctx, m)
}

func (s *Service) Update(ctx context.Context, id int64, url string, expectedStatus int, bodyContains string, intervalSeconds int, userAgent string) (domain.Monitor, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Monitor{}, err
	}

	validated, err := domain.NewMonitor(url, expectedStatus, bodyContains, intervalSeconds, userAgent)
	if err != nil {
		return domain.Monitor{}, err
	}

	updated := domain.Monitor{
		ID:              existing.ID,
		URL:             validated.URL,
		ExpectedStatus:  validated.ExpectedStatus,
		BodyContains:    validated.BodyContains,
		IntervalSeconds: validated.IntervalSeconds,
		UserAgent:       validated.UserAgent,
		CreatedAt:       existing.CreatedAt,
	}

	return s.repo.Update(ctx, updated)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
