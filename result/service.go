package result

import (
	"context"

	"github.com/sko/go-http-monitor/domain"
)

type Page struct {
	Results []domain.CheckResult `json:"results"`
	Total   int64                `json:"total"`
	Limit   int                  `json:"limit"`
	Offset  int                  `json:"offset"`
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Latest(ctx context.Context, monitorID int64) (domain.CheckResult, error) {
	return s.repo.FindLatestByMonitorID(ctx, monitorID)
}

func (s *Service) History(ctx context.Context, monitorID int64, limit, offset int) (Page, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	total, err := s.repo.CountByMonitorID(ctx, monitorID)
	if err != nil {
		return Page{}, err
	}

	results, err := s.repo.FindByMonitorID(ctx, monitorID, limit, offset)
	if err != nil {
		return Page{}, err
	}
	if results == nil {
		results = []domain.CheckResult{}
	}

	return Page{
		Results: results,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}
