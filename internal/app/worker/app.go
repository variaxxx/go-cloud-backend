package app_worker

import (
	app_metrics "cloud/internal/app/metrics"
	core_logger "cloud/internal/core/logger"
	"cloud/internal/features/file"
	infra_postgres "cloud/internal/infra/postgres"
	obs_prometheus "cloud/internal/observability/prometheus"
	"context"
	"fmt"
	"sync"
)

type runner interface {
	Run(ctx context.Context) error
	Close() error
}

type App struct {
	logger  *core_logger.Logger
	dbPool  infra_postgres.Pool
	runners []runner
}

func New(
	ctx context.Context,
) (*App, error) {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load WORKER config: %w", err)
	}

	dbConfig, err := infra_postgres.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load DB config: %w", err)
	}

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	db, err := infra_postgres.NewConnectionPool(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("db connection pool create: %w", err)
	}

	observability := obs_prometheus.RegisterWorker()

	fileUploadedWorker, err := file.NewUploadedWorker(file.WorkerDeps{
		DB:            db,
		Observability: observability,
	})
	if err != nil {
		return nil, fmt.Errorf("file uploaded worker init: %w", err)
	}

	metricsServer := app_metrics.NewServer(
		config.MetricsAddr,
		config.MetricsShutdownTimeout,
		logger,
		observability.Handler(),
	)

	return &App{
		logger: logger,
		dbPool: db,
		runners: []runner{
			fileUploadedWorker,
			metricsServer,
		},
	}, nil
}

func (a *App) Run(
	ctx context.Context,
) error {
	ctx = core_logger.WithContext(ctx, a.logger)
	a.logger.Info("Worker app started")

	errCh := make(chan error, len(a.runners))
	doneCh := make(chan struct{})
	var wg sync.WaitGroup

	for _, currentRunner := range a.runners {
		wg.Add(1)

		go func(currentRunner runner) {
			defer wg.Done()

			if err := currentRunner.Run(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(currentRunner)
	}

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("worker run: %w", err)
	case <-ctx.Done():
		<-doneCh
		return nil
	case <-doneCh:
		return nil
	}
}

func (a *App) Close() {
	a.logger.Info("Shutting down...")

	for _, currentRunner := range a.runners {
		_ = currentRunner.Close()
	}

	a.dbPool.Close()
	a.logger.Close()
}
