package auth_http

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_response "cloud/internal/core/transport/http/response"
	auth_domain "cloud/internal/features/auth/domain"
	"context"
	"net/http"
	"strings"
)

const userIDContextKey string = "user_id"

func Auth(
	tokenManager auth_domain.TokenManager,
) core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rh := core_http_response.NewHTTPResponseHandler(log, w)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				rh.ErrorResponse(
					core_errors.ErrUnauthorized,
					"Unauthorized",
				)
				return
			}

			prefix := "Bearer "
			if !strings.HasPrefix(authHeader, prefix) {
				rh.ErrorResponse(
					core_errors.ErrUnauthorized,
					"Bearer token required",
				)
				return
			}

			token := strings.TrimPrefix(authHeader, prefix)

			userID, err := tokenManager.Parse(token)
			if err != nil {
				rh.ErrorResponse(
					core_errors.ErrUnauthorized,
					"Invalid access token",
				)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserIDFromContext(
	ctx context.Context,
) (int64, bool) {
	id, ok := ctx.Value(userIDContextKey).(int64)
	return id, ok
}
