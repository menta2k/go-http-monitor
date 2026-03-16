package monitor

import (
	"context"

	"github.com/sko/go-http-monitor/domain"
)

type Repository interface {
	FindAll(ctx context.Context) ([]domain.Monitor, error)
	FindByID(ctx context.Context, id int64) (domain.Monitor, error)
	Create(ctx context.Context, m domain.Monitor) (domain.Monitor, error)
	Update(ctx context.Context, m domain.Monitor) (domain.Monitor, error)
	Delete(ctx context.Context, id int64) error
}
