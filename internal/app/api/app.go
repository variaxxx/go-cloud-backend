package app_api

import (
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	"cloud/internal/features/auth"
	"cloud/internal/features/file"
	"cloud/internal/features/hello"
	infra_postgres "cloud/internal/infra/postgres"
	"context"
	"fmt"
)

type App struct {
	logger     *core_logger.Logger
	httpServer *core_http_server.HTTPServer
	dbPool     infra_postgres.Pool
}

func New(
	ctx context.Context,
) (*App, error) {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	httpConfig, err := core_http_server.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load HTTP config: %w", err)
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

	apiRouter := core_http_server.NewAPIRouter("")
	hello.RegisterHTTP(apiRouter)

	authRouter := core_http_server.NewAPIRouter("/auth")
	if err := auth.Register(auth.Deps{
		Router: authRouter,
		DB:     db,
	}); err != nil {
		return nil, fmt.Errorf("auth module init: %w", err)
	}

	fileRouter := core_http_server.NewAPIRouter("/files")
	if err := file.Register(file.Deps{
		Router: fileRouter,
		DB:     db,
	}); err != nil {
		return nil, fmt.Errorf("file module init: %w", err)
	}

	middlewares := []core_http_middleware.Middleware{
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	}

	httpServer := core_http_server.NewHTTPServer(
		httpConfig,
		*logger,
		middlewares...,
	)
	httpServer.RegisterAPIRouters(
		*apiRouter,
		*authRouter,
		*fileRouter,
	)

	return &App{
		logger:     logger,
		httpServer: httpServer,
		dbPool:     db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}

func (a *App) Close() {
	a.logger.Close()
	a.dbPool.Close()
}
