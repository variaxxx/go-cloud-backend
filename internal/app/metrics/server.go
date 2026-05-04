package app_metrics

import (
	core_logger "cloud/internal/core/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	logger          *core_logger.Logger
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(
	addr string,
	shutdownTimeout time.Duration,
	logger *core_logger.Logger,
	handler http.Handler,
) *Server {
	return &Server{
		logger: logger,
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *Server) Run(
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

func (s *Server) Close() error {
	if err := s.server.Close(); errors.Is(err, http.ErrServerClosed) {
		return nil
	} else if err != nil {
		return fmt.Errorf("close metrics server: %w", err)
	}

	return nil
}
