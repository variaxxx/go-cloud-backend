package core_http_auth

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_response "cloud/internal/core/transport/http/response"
	"context"
	"net/http"
	"strings"
)

type TokenParser interface {
	Parse(tokenString string) (int64, error)
}

type userIDContextKey struct{}

func Middleware(
	tokenParser TokenParser,
) core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rh := core_http_response.NewHTTPResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "Unauthorized")
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "Bearer token required")
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)
			userID, err := tokenParser.Parse(token)
			if err != nil {
				rh.ErrorResponse(core_errors.ErrUnauthorized, "Invalid access token")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(
	ctx context.Context,
) (int64, bool) {
	id, ok := ctx.Value(userIDContextKey{}).(int64)
	return id, ok
}
