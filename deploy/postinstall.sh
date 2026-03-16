#!/bin/bash
set -e

# Create system user if not exists
if ! id -u go-http-monitor &>/dev/null; then
    useradd --system --no-create-home --shell /usr/sbin/nologin go-http-monitor
fi

# Create data directory
mkdir -p /var/lib/go-http-monitor
chown go-http-monitor:go-http-monitor /var/lib/go-http-monitor
chmod 750 /var/lib/go-http-monitor

# Create config directory and env file if not exists
mkdir -p /etc/go-http-monitor
if [ ! -f /etc/go-http-monitor/env ]; then
    cp /usr/share/go-http-monitor/env.example /etc/go-http-monitor/env
    chmod 640 /etc/go-http-monitor/env
    chown root:go-http-monitor /etc/go-http-monitor/env
    echo "Created /etc/go-http-monitor/env — edit it before starting the service"
fi

# Reload systemd and enable service
systemctl daemon-reload
systemctl enable go-http-monitor.service

echo ""
echo "go-http-monitor installed successfully."
echo ""
echo "  1. Edit /etc/go-http-monitor/env (set JWT_SECRET and ADMIN_PASSWORD)"
echo "  2. Start: systemctl start go-http-monitor"
echo "  3. Status: systemctl status go-http-monitor"
echo "  4. Logs: journalctl -u go-http-monitor -f"
echo ""
