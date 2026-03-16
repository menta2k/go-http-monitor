# go-http-monitor

A self-contained HTTP URL monitoring service with a REST API, embedded web UI, and time-series analytics. Ships as a single binary.

## Features

- **URL Monitoring** — Configure monitors with expected status code, body string match, custom User-Agent, and check interval
- **Background Workers** — Per-monitor goroutines check URLs on schedule with graceful shutdown
- **Notifications** — Email (SMTP) and Slack webhook alerts on state transitions only (no flood)
  - OK → Fail: sends alert
  - Fail → Fail: silent
  - Fail → OK: sends recovery
- **Time-Series Analytics** — FrostDB columnar storage for response time, uptime, and status code distribution
- **Charts** — ECharts dashboard with response time, health status, and HTTP status code trends
- **JWT Authentication** — Public login endpoint, all other routes protected
- **Dark/Light Theme** — Persisted to localStorage, defaults to OS preference
- **Housekeeper** — Automatic purge of old SQLite check results on configurable schedule
- **Single Binary** — Vue.js frontend embedded via `go:embed`, no separate web server needed

## Quick Start

```bash
# Build
make build

# Run
JWT_SECRET=your-secret ADMIN_PASSWORD=your-pass ./bin/go-http-monitor

# Open http://localhost:8080
```

## Configuration

All configuration is via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `DB_PATH` | `./monitor.db` | SQLite database path |
| `TSDB_PATH` | `./tsdb-data` | FrostDB storage directory |
| `JWT_SECRET` | *(required)* | HMAC signing key for JWT tokens |
| `ADMIN_PASSWORD` | *(required)* | Admin user password |
| `ADMIN_USERNAME` | `admin` | Admin username |
| `JWT_TOKEN_TTL_HOURS` | `24` | JWT token expiry in hours |
| `HTTP_CLIENT_TIMEOUT` | `30` | HTTP client timeout in seconds |
| `HOUSEKEEP_INTERVAL_MIN` | `60` | Housekeeper run interval in minutes |
| `HOUSEKEEP_RETENTION_DAYS` | `30` | Keep SQLite check results for N days |
| `SMTP_HOST` | *(empty)* | SMTP server (empty = email disabled) |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_FROM` | *(empty)* | Sender email address |
| `SMTP_USERNAME` | *(empty)* | SMTP auth username |
| `SMTP_PASSWORD` | *(empty)* | SMTP auth password |

## API Endpoints

### Authentication
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/login` | Login, returns JWT token |

### Monitors
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/monitors` | List all monitors |
| POST | `/api/monitors` | Create a monitor |
| GET | `/api/monitors/{id}` | Get monitor details |
| PUT | `/api/monitors/{id}` | Update a monitor |
| DELETE | `/api/monitors/{id}` | Delete a monitor |

### Check Results
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/monitors/{id}/status` | Latest check result |
| GET | `/api/monitors/{id}/history` | Paginated history (`?limit=20&offset=0`) |

### Notifications
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/monitors/{id}/notifications` | List notifications for a monitor |
| POST | `/api/monitors/{id}/notifications` | Create notification (email or slack) |
| GET | `/api/notifications/{nid}` | Get notification details |
| PUT | `/api/notifications/{nid}` | Update notification |
| DELETE | `/api/notifications/{nid}` | Delete notification |

### Analytics (FrostDB)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/monitors/{id}/stats` | Summary stats (`?period=1h\|6h\|24h\|7d\|30d`) |
| GET | `/api/monitors/{id}/timeline` | Time-series buckets (`?period=1h&buckets=60`) |
| GET | `/api/monitors/{id}/status-codes` | Status code distribution |
| GET | `/api/monitors/{id}/status-code-timeline` | Status codes over time |

## API Examples

```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"yourpass"}' | jq -r '.data.token')

# Create monitor
curl -X POST http://localhost:8080/api/monitors \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://example.com",
    "expected_status": 200,
    "body_contains": "Example Domain",
    "interval_seconds": 30,
    "user_agent": "MyBot/1.0"
  }'

# Add Slack notification
curl -X POST http://localhost:8080/api/monitors/1/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{
    "type": "slack",
    "target": "https://hooks.slack.com/services/T.../B.../xxx",
    "enabled": true
  }'

# Get stats
curl http://localhost:8080/api/monitors/1/stats?period=24h \
  -H "Authorization: Bearer $TOKEN"
```

## Architecture

```
domain/          Domain types, validation, errors
config/          Environment-based configuration
database/        SQLite connection, migrations
auth/            JWT authentication (service, handler, middleware)
monitor/         Monitor CRUD (repository, service, handler, routes)
result/          Check result storage and history API
notification/    Notification config CRUD (email, slack)
notifier/        Alert dispatcher with state-transition detection
checker/         HTTP checker, per-monitor workers, scheduler
tsdb/            FrostDB time-series storage
stats/           Analytics service and API (summary, timeline, status codes)
housekeeper/     Periodic SQLite cleanup
response/        JSON response helpers
frontend/        Vue 3 + Tailwind + ECharts SPA
```

**Storage:**
- **SQLite** — Monitors, notifications, check result history (CRUD)
- **FrostDB** — Columnar time-series for analytics (response time, uptime, status codes)

## Installation

### From Release (deb/rpm)

```bash
# Download from GitHub Releases
sudo dpkg -i go-http-monitor_1.0.0_amd64.deb

# Configure
sudo vim /etc/go-http-monitor/env

# Start
sudo systemctl start go-http-monitor
sudo systemctl status go-http-monitor

# Logs
journalctl -u go-http-monitor -f
```

### From Source

```bash
make build                    # Build frontend + backend
sudo make install             # Install binary + systemd service
sudo vim /etc/go-http-monitor/env
sudo systemctl start go-http-monitor
```

### Uninstall

```bash
sudo make uninstall
```

## Development

```bash
# Full build (frontend + backend)
make build

# Dev mode (Go backend + Vite HMR)
make dev

# Run tests
make test

# Run linter
make lint
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build frontend + Go binary into `bin/` |
| `make frontend` | Build Vue SPA only |
| `make backend` | Compile Go binary only |
| `make run` | Build and run |
| `make dev` | Dev mode with HMR |
| `make test` | Run Go tests |
| `make lint` | Run go vet |
| `make install` | Install binary + systemd service |
| `make uninstall` | Remove service + binary |
| `make clean` | Remove all build artifacts |

## Tech Stack

- **Backend:** Go, SQLite (modernc.org/sqlite), FrostDB, JWT (golang-jwt)
- **Frontend:** Vue 3, Pinia, Vue Router, Tailwind CSS, ECharts
- **Build:** Vite, go:embed, nFPM (deb/rpm packaging)
- **CI:** GitHub Actions (build + release on tag)

## License

MIT
