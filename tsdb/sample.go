package tsdb

// Sample is the FrostDB row schema for check results as time-series data.
// FrostDB uses struct tags to define the columnar schema.
type Sample struct {
	Timestamp      int64 `frostdb:"timestamp,asc(0)"`
	MonitorID      int64 `frostdb:"monitor_id,asc(1)"`
	StatusCode     int64 `frostdb:"status_code"`
	ResponseTimeMs int64 `frostdb:"response_time_ms"`
	Healthy        int64 `frostdb:"healthy"`      // 1 = healthy, 0 = failing
	BodyMatched    int64 `frostdb:"body_matched"` // 1 = matched, 0 = not matched, -1 = N/A
	HasError       int64 `frostdb:"has_error"`    // 1 = error, 0 = no error
}
