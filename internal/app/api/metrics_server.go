package app_api

import (
	core_logger "cloud/internal/core/logger"
	obs_prometheus "cloud/internal/observability/prometheus"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type metricsServer struct {
	logger          *core_logger.Logger
	server          *http.Server
	shutdownTimeout time.Duration
}

func newMetricsServer(
	config Config,
	logger *core_logger.Logger,
	obs *obs_prometheus.Observability,
) *metricsServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", obs.Handler())

	return &metricsServer{
		logger: logger,
		server: &http.Server{
			Addr:    config.MetricsAddr,
			Handler: mux,
		},
		shutdownTimeout: config.MetricsShutdownTimeout,
	}
}

func (s *metricsServer) Run(
	ctx context.Context,
) error {
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		s.logger.Info("Metrics server started", zap.String("addr", s.server.Addr))

		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err, ok := <-errCh:
		if !ok {
			return nil
		}

		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.shutdownTimeout,
		)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			_ = s.server.Close()
			return fmt.Errorf("shutdown metrics server: %w", err)
		}

		if err, ok := <-errCh; ok && err != nil {
			return fmt.Errorf("metrics serve after shutdown: %w", err)
		}
	}

	return nil
}
