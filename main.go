package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/sko/go-http-monitor/auth"
	"github.com/sko/go-http-monitor/checker"
	"github.com/sko/go-http-monitor/config"
	"github.com/sko/go-http-monitor/database"
	"github.com/sko/go-http-monitor/domain"
	"github.com/sko/go-http-monitor/monitor"
	"github.com/sko/go-http-monitor/notification"
	"github.com/sko/go-http-monitor/notifier"
	"github.com/sko/go-http-monitor/result"
)

func main() {
	cfg := config.Load()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if cfg.AdminPassword == "" {
		log.Fatal("ADMIN_PASSWORD environment variable is required")
	}

	db, err := database.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// Auth
	users := map[string]string{
		cfg.AdminUsername: cfg.AdminPassword,
	}
	authSvc := auth.NewService(cfg.JWTSecret, time.Duration(cfg.JWTTokenTTLHours)*time.Hour, users)
	authHandler := auth.NewHandler(authSvc)

	// Repositories
	monitorRepo := monitor.NewSQLiteRepository(db)
	resultRepo := result.NewSQLiteRepository(db)
	notifRepo := notification.NewSQLiteRepository(db)

	// Services
	monitorSvc := monitor.NewService(monitorRepo)
	resultSvc := result.NewService(resultRepo)
	notifSvc := notification.NewService(notifRepo)

	// Notifier senders
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.HTTPClientTimeout) * time.Second,
	}

	senders := map[domain.NotificationType]notifier.Sender{
		domain.NotificationSlack: notifier.NewSlackSender(httpClient),
	}
	if cfg.SMTPHost != "" {
		senders[domain.NotificationEmail] = notifier.NewEmailSender(
			cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom, cfg.SMTPUsername, cfg.SMTPPassword,
		)
		log.Printf("email notifications enabled (SMTP: %s:%d)", cfg.SMTPHost, cfg.SMTPPort)
	} else {
		log.Println("email notifications disabled (SMTP_HOST not set)")
	}

	alerter := notifier.New(notifRepo, resultRepo, senders)

	// Scheduler
	scheduler := checker.NewScheduler(httpClient, resultRepo, alerter.Notify)

	syncScheduler := func() {
		monitors, err := monitorSvc.List(context.Background())
		if err != nil {
			log.Printf("failed to sync scheduler: %v", err)
			return
		}
		scheduler.Sync(monitors)
	}

	// Handlers
	monitorHandler := monitor.NewHandler(monitorSvc, syncScheduler)
	resultHandler := result.NewHandler(resultSvc)
	notifHandler := notification.NewHandler(notifSvc)

	// API mux
	apiMux := http.NewServeMux()
	auth.RegisterRoutes(apiMux, authHandler)
	monitor.RegisterRoutes(apiMux, monitorHandler)
	result.RegisterRoutes(apiMux, resultHandler)
	notification.RegisterRoutes(apiMux, notifHandler)

	// JWT middleware
	jwtMiddleware := auth.RequireAuth(authSvc)

	// Root mux
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/auth/") {
			apiMux.ServeHTTP(w, r)
			return
		}
		jwtMiddleware(apiMux).ServeHTTP(w, r)
	})

	mux.Handle("/", frontendHandler())

	// Initial sync
	syncScheduler()

	// Server
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	scheduler.StopAll()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
