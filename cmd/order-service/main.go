package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"e-commerce/internal/config"
	"e-commerce/internal/domain/order"
	"e-commerce/internal/handler"
	"e-commerce/internal/infrastructure/db"
	"e-commerce/internal/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)
	log.Info("service starting",
		"service", cfg.ServiceName,
		"addr", cfg.ServerAddress,
		"log_level", cfg.LogLevel,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DatabaseDSN, cfg.DBMaxConns, cfg.DBMinConns, cfg.DBConnMaxIdleTime)
	if err != nil {
		log.Error("failed to init db pool", "error", err)
		os.Exit(1)
	}

	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	orderRepo := db.NewOrderRepository(pool)
	svc := order.NewService(orderRepo)
	orderHandler := handler.NewOrderHandler(svc, log)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
	mux.HandleFunc("POST /orders", orderHandler.Create)
	mux.HandleFunc("GET /orders/{id}", orderHandler.GetByID)
	mux.HandleFunc("GET /orders", orderHandler.List)

	srv := &http.Server{
		Addr:         cfg.ServerAddress,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		log.Info("signal received, starting shutdown")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error("server shutdown failed", "error", err)
		}

		cancel()
	}()

	log.Info("server listening", "addr", cfg.ServerAddress)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("server failed", "error", err)
	}

	<-ctx.Done()
	log.Info("service stopped")
}
