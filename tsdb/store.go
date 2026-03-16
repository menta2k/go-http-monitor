package tsdb

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/polarsignals/frostdb"

	"github.com/sko/go-http-monitor/domain"
)

type Store struct {
	columnStore *frostdb.ColumnStore
	db          *frostdb.DB
	table       *frostdb.GenericTable[Sample]
	mu          sync.Mutex
}

func Open(storagePath string) (*Store, error) {
	opts := []frostdb.Option{
		frostdb.WithActiveMemorySize(64 * 1024 * 1024), // 64MB
	}
	if storagePath != "" {
		opts = append(opts, frostdb.WithStoragePath(storagePath))
	}

	cs, err := frostdb.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create frostdb column store: %w", err)
	}

	ctx := context.Background()
	db, err := cs.DB(ctx, "monitor_metrics")
	if err != nil {
		cs.Close()
		return nil, fmt.Errorf("create frostdb database: %w", err)
	}

	table, err := frostdb.NewGenericTable[Sample](
		db, "check_samples", memory.NewGoAllocator(),
	)
	if err != nil {
		cs.Close()
		return nil, fmt.Errorf("create frostdb table: %w", err)
	}

	log.Printf("frostdb opened (path: %s)", storagePath)
	return &Store{
		columnStore: cs,
		db:          db,
		table:       table,
	}, nil
}

func (s *Store) Write(ctx context.Context, m domain.Monitor, cr domain.CheckResult) error {
	sample := toSample(m, cr)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := s.table.Write(ctx, sample); err != nil {
		return fmt.Errorf("write sample: %w", err)
	}
	return nil
}

func (s *Store) Table() *frostdb.Table {
	return s.table.Table
}

func (s *Store) GenericTable() *frostdb.GenericTable[Sample] {
	return s.table
}

func (s *Store) Close() error {
	return s.columnStore.Close()
}

func toSample(m domain.Monitor, cr domain.CheckResult) Sample {
	healthy := int64(1)
	if cr.Error != "" || cr.StatusCode != m.ExpectedStatus {
		healthy = 0
	}
	if cr.BodyMatched != nil && !*cr.BodyMatched {
		healthy = 0
	}

	bodyMatched := int64(-1) // N/A
	if cr.BodyMatched != nil {
		if *cr.BodyMatched {
			bodyMatched = 1
		} else {
			bodyMatched = 0
		}
	}

	hasError := int64(0)
	if cr.Error != "" {
		hasError = 1
	}

	return Sample{
		Timestamp:      cr.CheckedAt.UnixMilli(),
		MonitorID:      m.ID,
		StatusCode:     int64(cr.StatusCode),
		ResponseTimeMs: cr.ResponseTimeMs,
		Healthy:        healthy,
		BodyMatched:    bodyMatched,
		HasError:       hasError,
	}
}
