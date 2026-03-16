package checker

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sko/go-http-monitor/domain"
)

const maxBodySize = 1 << 20 // 1MB

func Check(ctx context.Context, client *http.Client, m domain.Monitor) domain.CheckResult {
	now := time.Now().UTC()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return domain.CheckResult{
			MonitorID: m.ID,
			Error:     err.Error(),
			CheckedAt: now,
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return domain.CheckResult{
			MonitorID:      m.ID,
			ResponseTimeMs: elapsed,
			Error:          err.Error(),
			CheckedAt:      now,
		}
	}
	defer resp.Body.Close()

	result := domain.CheckResult{
		MonitorID:      m.ID,
		StatusCode:     resp.StatusCode,
		ResponseTimeMs: elapsed,
		CheckedAt:      now,
	}

	if m.BodyContains != "" {
		body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodySize))
		if err != nil {
			result.Error = "failed to read body: " + err.Error()
			return result
		}
		matched := strings.Contains(string(body), m.BodyContains)
		result.BodyMatched = &matched
	}

	return result
}
