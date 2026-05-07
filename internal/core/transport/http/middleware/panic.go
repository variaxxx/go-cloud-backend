package core_http_middleware

import (
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	"net/http"
)

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"Unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
