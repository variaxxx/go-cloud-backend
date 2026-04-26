package core_http_response

import (
	core_errors "cloud/internal/core/errors"
	core_logger "cloud/internal/core/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("Write HTTP JSON response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	status := http.StatusInternalServerError
	err := fmt.Errorf("Unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.rw.WriteHeader(status)

	h.errorResponse(status, msg)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, msg)
}

func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	msg string,
) {
	response := map[string]string{
		"statusCode": strconv.Itoa(statusCode),
		"message":    msg,
	}

	h.JSONResponse(response, statusCode)
}
