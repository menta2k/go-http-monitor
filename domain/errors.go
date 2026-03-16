package domain

import "errors"

var (
	ErrMonitorNotFound  = errors.New("monitor not found")
	ErrInvalidURL       = errors.New("invalid URL")
	ErrInvalidStatusCode = errors.New("invalid status code")
	ErrInvalidInterval  = errors.New("invalid interval")
)
