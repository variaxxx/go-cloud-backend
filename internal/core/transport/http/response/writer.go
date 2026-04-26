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

func (rw *HTTPResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("Tried to get uninitialized status code")
	}
	return rw.statusCode
}
