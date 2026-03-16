package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/polarsignals/frostdb/query/logicalplan"

	"github.com/sko/go-http-monitor/tsdb"
)

type Summary struct {
	MonitorID     int64   `json:"monitor_id"`
	Period        string  `json:"period"`
	TotalChecks   int64   `json:"total_checks"`
	HealthyChecks int64   `json:"healthy_checks"`
	FailedChecks  int64   `json:"failed_checks"`
	ErrorChecks   int64   `json:"error_checks"`
	UptimePct     float64 `json:"uptime_pct"`
	AvgResponseMs float64 `json:"avg_response_ms"`
	MinResponseMs int64   `json:"min_response_ms"`
	MaxResponseMs int64   `json:"max_response_ms"`
	P95ResponseMs int64   `json:"p95_response_ms"`
}

type StatusCodeBucket struct {
	Code  string `json:"code"`
	Count int64  `json:"count"`
}

type StatusCodeTimePoint struct {
	Timestamp int64            `json:"timestamp"`
	Codes     map[string]int64 `json:"codes"`
}

type TimePoint struct {
	Timestamp     int64   `json:"timestamp"`
	AvgResponseMs float64 `json:"avg_response_ms"`
	Healthy       int64   `json:"healthy"`
	Total         int64   `json:"total"`
}

type Service struct {
	store *tsdb.Store
}

func NewService(store *tsdb.Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetSummary(ctx context.Context, monitorID int64, period string) (Summary, error) {
	since := periodToTime(period)
	sinceMs := since.UnixMilli()

	var samples []tsdb.Sample
	err := scanSamples(ctx, s.store, monitorID, sinceMs, func(sp tsdb.Sample) {
		samples = append(samples, sp)
	})
	if err != nil {
		return Summary{}, err
	}

	if len(samples) == 0 {
		return Summary{MonitorID: monitorID, Period: period}, nil
	}

	total := int64(len(samples))
	var healthy, errors int64
	var sumResp, minResp, maxResp int64
	minResp = samples[0].ResponseTimeMs

	for _, sp := range samples {
		healthy += sp.Healthy
		errors += sp.HasError
		sumResp += sp.ResponseTimeMs
		if sp.ResponseTimeMs < minResp {
			minResp = sp.ResponseTimeMs
		}
		if sp.ResponseTimeMs > maxResp {
			maxResp = sp.ResponseTimeMs
		}
	}

	return Summary{
		MonitorID:     monitorID,
		Period:        period,
		TotalChecks:   total,
		HealthyChecks: healthy,
		FailedChecks:  total - healthy,
		ErrorChecks:   errors,
		UptimePct:     float64(healthy) / float64(total) * 100,
		AvgResponseMs: float64(sumResp) / float64(total),
		MinResponseMs: minResp,
		MaxResponseMs: maxResp,
		P95ResponseMs: computeP95(samples),
	}, nil
}

func (s *Service) GetTimeline(ctx context.Context, monitorID int64, period string, buckets int) ([]TimePoint, error) {
	since := periodToTime(period)
	sinceMs := since.UnixMilli()
	nowMs := time.Now().UnixMilli()

	if buckets <= 0 || buckets > 200 {
		buckets = 60
	}

	var samples []tsdb.Sample
	err := scanSamples(ctx, s.store, monitorID, sinceMs, func(sp tsdb.Sample) {
		samples = append(samples, sp)
	})
	if err != nil {
		return nil, err
	}

	if len(samples) == 0 {
		return []TimePoint{}, nil
	}

	bucketSize := (nowMs - sinceMs) / int64(buckets)
	if bucketSize <= 0 {
		bucketSize = 1
	}

	points := make([]TimePoint, buckets)
	for i := range points {
		points[i].Timestamp = sinceMs + int64(i)*bucketSize + bucketSize/2
	}

	for _, sp := range samples {
		idx := int((sp.Timestamp - sinceMs) / bucketSize)
		if idx < 0 {
			idx = 0
		}
		if idx >= buckets {
			idx = buckets - 1
		}
		points[idx].Total++
		points[idx].Healthy += sp.Healthy
		points[idx].AvgResponseMs += float64(sp.ResponseTimeMs)
	}

	for i := range points {
		if points[i].Total > 0 {
			points[i].AvgResponseMs = points[i].AvgResponseMs / float64(points[i].Total)
		}
	}

	return points, nil
}

func (s *Service) GetStatusCodes(ctx context.Context, monitorID int64, period string) ([]StatusCodeBucket, error) {
	since := periodToTime(period)
	sinceMs := since.UnixMilli()

	counts := map[string]int64{}
	err := scanSamples(ctx, s.store, monitorID, sinceMs, func(sp tsdb.Sample) {
		counts[statusCodeBucket(sp)]++
	})
	if err != nil {
		return nil, err
	}

	// Return in a stable order
	order := []string{"2xx", "3xx", "4xx", "5xx", "Error"}
	var result []StatusCodeBucket
	for _, code := range order {
		if c, ok := counts[code]; ok {
			result = append(result, StatusCodeBucket{Code: code, Count: c})
		}
	}
	if len(result) == 0 {
		return []StatusCodeBucket{}, nil
	}
	return result, nil
}

func (s *Service) GetStatusCodeTimeline(ctx context.Context, monitorID int64, period string, buckets int) ([]StatusCodeTimePoint, error) {
	since := periodToTime(period)
	sinceMs := since.UnixMilli()
	nowMs := time.Now().UnixMilli()

	if buckets <= 0 || buckets > 200 {
		buckets = 60
	}

	var samples []tsdb.Sample
	err := scanSamples(ctx, s.store, monitorID, sinceMs, func(sp tsdb.Sample) {
		samples = append(samples, sp)
	})
	if err != nil {
		return nil, err
	}

	bucketSize := (nowMs - sinceMs) / int64(buckets)
	if bucketSize <= 0 {
		bucketSize = 1
	}

	points := make([]StatusCodeTimePoint, buckets)
	for i := range points {
		points[i].Timestamp = sinceMs + int64(i)*bucketSize + bucketSize/2
		points[i].Codes = map[string]int64{}
	}

	for _, sp := range samples {
		idx := int((sp.Timestamp - sinceMs) / bucketSize)
		if idx < 0 {
			idx = 0
		}
		if idx >= buckets {
			idx = buckets - 1
		}

		bucket := statusCodeBucket(sp)
		points[idx].Codes[bucket]++
	}

	return points, nil
}

func statusCodeBucket(sp tsdb.Sample) string {
	if sp.HasError == 1 {
		return "Error"
	}
	switch {
	case sp.StatusCode >= 200 && sp.StatusCode < 300:
		return "2xx"
	case sp.StatusCode >= 300 && sp.StatusCode < 400:
		return "3xx"
	case sp.StatusCode >= 400 && sp.StatusCode < 500:
		return "4xx"
	case sp.StatusCode >= 500:
		return "5xx"
	default:
		return "Error"
	}
}

func scanSamples(ctx context.Context, store *tsdb.Store, monitorID, sinceMs int64, fn func(tsdb.Sample)) error {
	table := store.Table()
	pool := memory.NewGoAllocator()

	return table.View(ctx, func(ctx context.Context, tx uint64) error {
		return table.Iterator(ctx, tx, pool, []logicalplan.Callback{
			func(_ context.Context, r arrow.Record) error {
				schema := r.Schema()
				monCol := r.Column(schema.FieldIndices("monitor_id")[0]).(*array.Int64)
				tsCol := r.Column(schema.FieldIndices("timestamp")[0]).(*array.Int64)
				scCol := r.Column(schema.FieldIndices("status_code")[0]).(*array.Int64)
				rtCol := r.Column(schema.FieldIndices("response_time_ms")[0]).(*array.Int64)
				hCol := r.Column(schema.FieldIndices("healthy")[0]).(*array.Int64)
				bmCol := r.Column(schema.FieldIndices("body_matched")[0]).(*array.Int64)
				heCol := r.Column(schema.FieldIndices("has_error")[0]).(*array.Int64)

				for i := 0; i < int(r.NumRows()); i++ {
					mid := monCol.Value(i)
					ts := tsCol.Value(i)
					if mid != monitorID || ts < sinceMs {
						continue
					}
					fn(tsdb.Sample{
						Timestamp:      ts,
						MonitorID:      mid,
						StatusCode:     scCol.Value(i),
						ResponseTimeMs: rtCol.Value(i),
						Healthy:        hCol.Value(i),
						BodyMatched:    bmCol.Value(i),
						HasError:       heCol.Value(i),
					})
				}
				return nil
			},
		})
	})
}

func computeP95(samples []tsdb.Sample) int64 {
	n := len(samples)
	if n == 0 {
		return 0
	}

	times := make([]int64, n)
	for i, s := range samples {
		times[i] = s.ResponseTimeMs
	}
	sortInt64s(times)

	idx := int(float64(n) * 0.95)
	if idx >= n {
		idx = n - 1
	}
	return times[idx]
}

func sortInt64s(a []int64) {
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && a[j] > key {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

func periodToTime(period string) time.Time {
	now := time.Now().UTC()
	switch period {
	case "1h":
		return now.Add(-1 * time.Hour)
	case "6h":
		return now.Add(-6 * time.Hour)
	case "24h", "1d":
		return now.Add(-24 * time.Hour)
	case "7d":
		return now.Add(-7 * 24 * time.Hour)
	case "30d":
		return now.Add(-30 * 24 * time.Hour)
	default:
		return now.Add(-24 * time.Hour)
	}
}

func ParsePeriod(s string) string {
	switch s {
	case "1h", "6h", "24h", "1d", "7d", "30d":
		return s
	default:
		return "24h"
	}
}

func ParseBuckets(n int) int {
	if n <= 0 || n > 200 {
		return 60
	}
	return n
}

func FormatDuration(period string) string {
	switch period {
	case "1h":
		return "Last 1 hour"
	case "6h":
		return "Last 6 hours"
	case "24h", "1d":
		return "Last 24 hours"
	case "7d":
		return "Last 7 days"
	case "30d":
		return "Last 30 days"
	default:
		return fmt.Sprintf("Last %s", period)
	}
}
