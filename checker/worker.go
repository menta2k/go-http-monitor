package checker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/result"
)

// AlertFunc is called after each check with the monitor and result.
type AlertFunc func(ctx context.Context, m domain.Monitor, cr domain.CheckResult)

func RunWorker(ctx context.Context, client *http.Client, m domain.Monitor, repo result.Repository, alert AlertFunc) {
	ticker := time.NewTicker(time.Duration(m.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	log.Printf("[worker] started monitoring %s (id=%d, every %ds)", m.URL, m.ID, m.IntervalSeconds)

	performCheck(ctx, client, m, repo, alert)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[worker] stopped monitoring %s (id=%d)", m.URL, m.ID)
			return
		case <-ticker.C:
			performCheck(ctx, client, m, repo, alert)
		}
	}
}

func performCheck(ctx context.Context, client *http.Client, m domain.Monitor, repo result.Repository, alert AlertFunc) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[worker] panic checking %s (id=%d): %v", m.URL, m.ID, r)
		}
	}()

	cr := Check(ctx, client, m)
	if _, err := repo.Create(ctx, cr); err != nil {
		log.Printf("[worker] failed to save result for %s (id=%d): %v", m.URL, m.ID, err)
		return
	}

	status := "OK"
	if cr.Error != "" {
		status = "ERROR: " + cr.Error
	} else if cr.StatusCode != m.ExpectedStatus {
		status = "STATUS_MISMATCH"
	} else if cr.BodyMatched != nil && !*cr.BodyMatched {
		status = "BODY_MISMATCH"
	}
	log.Printf("[worker] check %s (id=%d): %s (status=%d, time=%dms)", m.URL, m.ID, status, cr.StatusCode, cr.ResponseTimeMs)

	if alert != nil {
		alert(ctx, m, cr)
	}
}
