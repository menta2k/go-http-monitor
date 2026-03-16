package checker

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/result"
)

type Scheduler struct {
	client     *http.Client
	resultRepo result.Repository
	alert      AlertFunc
	mu         sync.Mutex
	workers    map[int64]context.CancelFunc
}

func NewScheduler(client *http.Client, resultRepo result.Repository, alert AlertFunc) *Scheduler {
	return &Scheduler{
		client:     client,
		resultRepo: resultRepo,
		alert:      alert,
		workers:    make(map[int64]context.CancelFunc),
	}
}

func (s *Scheduler) Sync(monitors []domain.Monitor) {
	s.mu.Lock()
	defer s.mu.Unlock()

	active := make(map[int64]bool, len(monitors))

	for _, m := range monitors {
		active[m.ID] = true
		if _, running := s.workers[m.ID]; !running {
			s.startLocked(m)
		}
	}

	for id, cancel := range s.workers {
		if !active[id] {
			cancel()
			delete(s.workers, id)
		}
	}
}

func (s *Scheduler) Start(m domain.Monitor) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, exists := s.workers[m.ID]; exists {
		cancel()
		delete(s.workers, m.ID)
	}

	s.startLocked(m)
}

func (s *Scheduler) Stop(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if cancel, exists := s.workers[id]; exists {
		cancel()
		delete(s.workers, id)
	}
}

func (s *Scheduler) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, cancel := range s.workers {
		cancel()
		delete(s.workers, id)
	}
	log.Println("[scheduler] all workers stopped")
}

func (s *Scheduler) startLocked(m domain.Monitor) {
	ctx, cancel := context.WithCancel(context.Background())
	s.workers[m.ID] = cancel
	go RunWorker(ctx, s.client, m, s.resultRepo, s.alert)
}
