package core_http_response

import "net/http"

const (
	StatusCodeUninitialized = -1
)

type HTTPResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewHTTPResponseWriter(
	rw http.ResponseWriter,
) *HTTPResponseWriter {
	return &HTTPResponseWriter{
		ResponseWriter: rw,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *HTTPResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *HTTPResponseWriter) Write(data []byte) (int, error) {
	if rw.statusCode == StatusCodeUninitialized {
		rw.statusCode = http.StatusOK
	}

	return rw.ResponseWriter.Write(data)
}

func (rw *HTTPResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}
	return rw.statusCode
}
