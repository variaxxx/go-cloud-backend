package app_api

import (
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_server "cloud/internal/core/transport/http/server"
	"cloud/internal/features/hello"
	"context"
	"fmt"
)

type App struct {
	logger     *core_logger.Logger
	httpServer *core_http_server.HTTPServer
}

func New() (*App, error) {
	loggerConfig, err := core_logger.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	httpConfig, err := core_http_server.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load HTTP config: %w", err)
	}

	logger, err := core_logger.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	apiRouter := core_http_server.NewAPIRouter("")
	hello.RegisterHTTP(apiRouter)

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
	httpServer.RegisterAPIRouters(*apiRouter)

	return &App{
		logger:     logger,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.httpServer.Run(ctx)
}

func (a *App) Close() {
	a.logger.Close()
}
