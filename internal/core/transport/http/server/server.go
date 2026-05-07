package core_http_server

import (
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	cfg    Config
	logger core_logger.Logger

	routes     []string
	middleware []core_http_middleware.Middleware
}

func NewHTTPServer(
	cfg Config,
	logger core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		cfg:        cfg,
		logger:     logger,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...APIRouter) {
	for _, router := range routers {
		s.mux.Handle(
			router.prefix+"/",
			http.StripPrefix(router.prefix, router),
		)

		for _, route := range router.routes {
			endpoint := fmt.Sprintf("%s %s%s", route.Method, router.prefix, route.Path)

			s.routes = append(
				s.routes,
				endpoint,
			)

			s.logger.Debug("Registered route", zap.String("route", endpoint))
		}
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Info("Starting HTTP server...")
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: mux,
	}

	listener, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return fmt.Errorf("HTTP listen: %w", err)
	}

	s.logger.Info("HTTP server started", zap.String("addr", s.cfg.Addr))

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)

		if err := server.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("HTTP listen and serve: %w", err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		s.logger.Info("Shutting server down...")

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("HTTP server shutdown: %w", err)
		}

		if err, ok := <-errCh; ok && err != nil {
			return fmt.Errorf("HTTP serve after shutdown: %w", err)
		}
	}

	return nil
}
