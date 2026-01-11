package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saisudhir14/fleet-event-stream/internal/handlers"
	"github.com/saisudhir14/fleet-event-stream/internal/metrics"
	"github.com/saisudhir14/fleet-event-stream/internal/processor"
)

// Constants for default configuration values
const (
	defaultPort        = "8080"
	defaultMetricsPort = "9090"
	shutdownTimeout    = 30 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	// Initialize components
	m := metrics.NewMetrics()
	proc := processor.NewEventProcessor(logger)
	h := handlers.NewHandler(proc, m, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ready", h.ReadyCheck)
	mux.HandleFunc("/api/v1/events", h.IngestEvent)
	mux.HandleFunc("/api/v1/stats", h.GetStats)
	port := getEnv("PORT", defaultPort)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// Metrics server
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsPort := getEnv("METRICS_PORT", defaultMetricsPort)
	metricsServer := &http.Server{

		Addr:    fmt.Sprintf(":%s", metricsPort),
		Handler: metricsMux,
	}
	// Start servers in separate goroutines
	go func() {
		logger.Info("starting metrics server", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("metrics server failed", "error", err)
		}
	}()
	// Start servers in separate goroutines
	go func() {
		logger.Info("starting API server", "port", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("API server failed", "error", err)
		}
	}()
	// Graceful shutdown on interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down servers...")

	// Creating a context with timeout for the shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown servers
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("API server shutdown failed", "error", err)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		logger.Error("metrics server shutdown failed", "error", err)
	}
	logger.Info("servers stopped gracefully")

}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is empty or not present, it returns the defaultValue.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
