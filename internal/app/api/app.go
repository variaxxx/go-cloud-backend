package app_api

import (
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	"cloud/internal/features/auth"
	"cloud/internal/features/file"
	"cloud/internal/features/folder"
	"cloud/internal/features/test"
	"cloud/internal/features/user"
	infra_postgres "cloud/internal/infra/postgres"
	obs_prometheus "cloud/internal/observability/prometheus"
	"context"
	"fmt"
	"sync"
)

type App struct {
	logger        *core_logger.Logger
	httpServer    *core_http_server.HTTPServer
	metricsServer *metricsServer
	dbPool        infra_postgres.Pool
}

type routerRegistration struct {
	prefix       string
	moduleName   string
	registerFunc func(router *core_http_server.APIRouter) error
}

func New(
	ctx context.Context,
) (*App, error) {
	logger, err := initLogger()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			logger.Close()
		}
	}()

	config, err := NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load API config: %w", err)
	}

	httpConfig, err := initHTTPConfig()
	if err != nil {
		return nil, err
	}

	db, err := initDB(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	routers, err := registerRouters(db)
	if err != nil {
		return nil, err
	}

	observability := obs_prometheus.Register()

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		*logger,
		buildMiddlewares(logger, observability.HTTP)...,
	)
	httpServer.RegisterAPIRouters(routers...)

	metricsServer := newMetricsServer(config, logger, observability)

	return &App{
		logger:        logger,
		httpServer:    httpServer,
		metricsServer: metricsServer,
		dbPool:        db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	doneCh := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := a.httpServer.Run(runCtx); err != nil {
			select {
			case errCh <- fmt.Errorf("run HTTP server: %w", err):
			default:
			}
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := a.metricsServer.Run(runCtx); err != nil {
			select {
			case errCh <- fmt.Errorf("run metrics server: %w", err):
			default:
			}
			cancel()
		}
	}()

	go func() {
		defer close(doneCh)
		wg.Wait()
	}()

	select {
	case err := <-errCh:
		cancel()
		<-doneCh
		return err
	case <-runCtx.Done():
		<-doneCh
		return nil
	case <-doneCh:
		select {
		case err := <-errCh:
			return err
		default:
			return nil
		}
	}
}

func (a *App) Close() {
	a.logger.Close()
	a.dbPool.Close()
}

func initHTTPConfig() (core_http_server.Config, error) {
	httpConfig, err := core_http_server.NewConfig()
	if err != nil {
		return core_http_server.Config{}, fmt.Errorf("load HTTP config: %w", err)
	}

	return httpConfig, nil
}

func initLogger() (*core_logger.Logger, error) {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	return logger, nil
}

func initDB(
	ctx context.Context,
) (*infra_postgres.ConnectionPool, error) {
	dbConfig, err := infra_postgres.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load DB config: %w", err)
	}

	db, err := infra_postgres.NewConnectionPool(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("db connection pool create: %w", err)
	}

	return db, nil
}

func buildMiddlewares(
	logger *core_logger.Logger,
	httpMetrics *obs_prometheus.HTTPMetrics,
) []core_http_middleware.Middleware {
	return []core_http_middleware.Middleware{
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		httpMetrics.Middleware(),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	}
}

func registerRouters(
	db infra_postgres.Pool,
) ([]core_http_server.APIRouter, error) {
	registrations := []routerRegistration{
		{
			prefix:     "",
			moduleName: "test",
			registerFunc: func(router *core_http_server.APIRouter) error {
				return test.Register(test.Deps{
					Router: router,
					DB:     db,
				})
			},
		},
		{
			prefix:     "/auth",
			moduleName: "auth",
			registerFunc: func(router *core_http_server.APIRouter) error {
				return auth.Register(auth.Deps{
					Router: router,
					DB:     db,
				})
			},
		},
		{
			prefix:     "/files",
			moduleName: "file",
			registerFunc: func(router *core_http_server.APIRouter) error {
				return file.Register(file.Deps{
					Router: router,
					DB:     db,
				})
			},
		},
		{
			prefix:     "/folders",
			moduleName: "folder",
			registerFunc: func(router *core_http_server.APIRouter) error {
				return folder.Register(folder.Deps{
					Router: router,
					DB:     db,
				})
			},
		},
		{
			prefix:     "/me",
			moduleName: "user",
			registerFunc: func(router *core_http_server.APIRouter) error {
				return user.Register(user.Deps{
					Router: router,
					DB:     db,
				})
			},
		},
	}

	routers := make([]core_http_server.APIRouter, 0, len(registrations))
	for _, registration := range registrations {
		router := core_http_server.NewAPIRouter(registration.prefix)
		if err := registration.registerFunc(router); err != nil {
			return nil, fmt.Errorf("%s module init: %w", registration.moduleName, err)
		}

		routers = append(routers, *router)
	}

	return routers, nil
}
