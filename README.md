# go-http-monitor

> Because refreshing a browser tab and whispering "please be up" is not a monitoring strategy.

A single-binary URL monitoring service that watches your websites so you don't have to. It pings your URLs, remembers everything in a time-series database, draws pretty charts, and screams at you on Slack or email when things go sideways.

One binary. No Docker required. No Kubernetes. No "let me just spin up 14 microservices." Just run it.

## What it does

- **Watches your URLs** and checks status codes, response times, and whether a specific string is in the response body
- **Sends you alerts** via Email or Slack -- but only when something actually *changes* (no 3am notification floods because a server is still down)
- **Draws charts** so you can pretend to understand what P95 means in your next standup
- **Cleans up after itself** with a built-in housekeeper that purges old data so your disk doesn't fill up in 2027
- **Looks good doing it** with a dark/light theme that respects your OS preference

## Getting started

```bash
make build

JWT_SECRET=something-secret ADMIN_PASSWORD=not-password123 ./bin/go-http-monitor
```

Open `http://localhost:8080`, log in, add a URL, and go grab a coffee. It's monitoring now.

## The dashboard

The monitor detail page gives you:

- **Uptime percentage** -- the number you paste into the SLA spreadsheet
- **Response time chart** -- a pretty line that hopefully stays flat and low
- **Health status bars** -- green means good, red means someone's getting paged
- **HTTP status code trends** -- see exactly when your API started returning 500s

All powered by [FrostDB](https://github.com/polarsignals/frostdb), a columnar time-series database. Because storing metrics in a regular SQL table is like using a spreadsheet as a database -- it works until it doesn't.

## Notifications that don't suck

Most monitors spam you with "DOWN DOWN DOWN DOWN" every 60 seconds. This one doesn't.

| What happens | What you get |
|---|---|
| Site goes down | One alert |
| Site is still down | Nothing (you already know) |
| Site comes back | One recovery message |
| Site is still up | Nothing (as it should be) |

Set up as many notification channels per monitor as you want -- mix email and Slack, send to different channels, whatever. Configure it all through the UI or API.

## Configuration

Everything is an environment variable. No YAML files were harmed in the making of this project.

| Variable | Default | What it does |
|----------|---------|-------------|
| `PORT` | `8080` | Where the web UI lives |
| `DB_PATH` | `./monitor.db` | SQLite file for monitors and history |
| `TSDB_PATH` | `./tsdb-data` | FrostDB directory for time-series data |
| `JWT_SECRET` | *(you must set this)* | Secret for auth tokens |
| `ADMIN_PASSWORD` | *(you must set this)* | Your login password |
| `ADMIN_USERNAME` | `admin` | Your login username |
| `HOUSEKEEP_INTERVAL_MIN` | `60` | How often to clean old data (minutes) |
| `HOUSEKEEP_RETENTION_DAYS` | `30` | How long to keep check history |
| `SMTP_HOST` | *(empty = no email)* | SMTP server for email alerts |
| `SMTP_PORT` | `587` | SMTP port |
| `SMTP_FROM` | | Sender address |
| `SMTP_USERNAME` | | SMTP login |
| `SMTP_PASSWORD` | | SMTP password |
| `HTTP_CLIENT_TIMEOUT` | `30` | Seconds before giving up on a check |
| `JWT_TOKEN_TTL_HOURS` | `24` | How long login sessions last |

## API

Full REST API behind JWT auth. Login at `POST /api/auth/login`, then use the token as `Authorization: Bearer <token>`.

```bash
# Get a token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"yourpass"}' | jq -r '.data.token')

# Add a monitor
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

# Add a Slack alert
curl -X POST http://localhost:8080/api/monitors/1/notifications \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"type":"slack","target":"https://hooks.slack.com/services/...","enabled":true}'

# Check the stats
curl http://localhost:8080/api/monitors/1/stats?period=24h \
  -H "Authorization: Bearer $TOKEN"
```

**Endpoints at a glance:**

| What | Endpoint |
|------|----------|
| Login | `POST /api/auth/login` |
| Monitors CRUD | `GET/POST /api/monitors`, `GET/PUT/DELETE /api/monitors/{id}` |
| Latest check | `GET /api/monitors/{id}/status` |
| Check history | `GET /api/monitors/{id}/history?limit=20&offset=0` |
| Notifications | `GET/POST /api/monitors/{id}/notifications`, `GET/PUT/DELETE /api/notifications/{nid}` |
| Stats summary | `GET /api/monitors/{id}/stats?period=24h` |
| Timeline | `GET /api/monitors/{id}/timeline?period=1h&buckets=60` |
| Status codes | `GET /api/monitors/{id}/status-codes?period=24h` |
| Status code timeline | `GET /api/monitors/{id}/status-code-timeline?period=1h&buckets=60` |

## Installing on a server

### The easy way (deb/rpm)

```bash
sudo dpkg -i go-http-monitor_1.0.0_amd64.deb
sudo vim /etc/go-http-monitor/env    # set JWT_SECRET and ADMIN_PASSWORD
sudo systemctl start go-http-monitor
```

It creates a systemd service, a dedicated user, and a config file. Like a proper grown-up application.

### The from-source way

```bash
make build
sudo make install
sudo vim /etc/go-http-monitor/env
sudo systemctl start go-http-monitor
```

### Check the logs

```bash
journalctl -u go-http-monitor -f
```

## Development

```bash
make build    # Build everything
make dev      # Go backend + Vite HMR (hot reload for frontend)
make test     # Run tests
make lint     # Run go vet
make clean    # Nuke build artifacts
```

## How it's built

- **Go** -- because it compiles to a single binary and goroutines are basically free
- **SQLite** -- for the boring CRUD stuff (monitors, notifications, check history)
- **FrostDB** -- for the cool time-series stuff (analytics, charts)
- **Vue 3 + Tailwind** -- because the frontend needs to look nice
- **ECharts** -- for the charts that make you feel like you're running a NASA control room
- **go:embed** -- the entire frontend is baked into the binary. One file to rule them all.

## License

MIT -- do whatever you want with it. If it breaks your production, that's between you and your SLA.
