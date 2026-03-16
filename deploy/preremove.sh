#!/bin/bash
set -e

systemctl stop go-http-monitor.service 2>/dev/null || true
systemctl disable go-http-monitor.service 2>/dev/null || true
