package obs_prometheus

import (
	core_http_middleware "cloud/internal/core/transport/http/middleware"
	core_http_response "cloud/internal/core/transport/http/response"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type HTTPMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "cloud",
				Name:      "http_requests_total",
				Help:      "Total count of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "cloud",
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
	}
}

func (m *HTTPMetrics) Middleware() core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			rw := core_http_response.NewHTTPResponseWriter(w)

			next.ServeHTTP(rw, r)

			statusCode := strconv.Itoa(rw.GetStatusCode())
			path := r.Pattern
			if path == "" {
				path = r.URL.Path
			}

			labels := prometheus.Labels{
				"method": r.Method,
				"path":   path,
				"status": statusCode,
			}

			m.RequestsTotal.With(labels).Inc()
			m.RequestDuration.With(labels).Observe(time.Since(startedAt).Seconds())
		})
	}
}
