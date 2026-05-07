package test_http

import (
	core_logger "cloud/internal/core/logger"
	core_http_response "cloud/internal/core/transport/http/response"
	"net/http"
)

func (h *Handler) GetHello(
	rw http.ResponseWriter,
	r *http.Request,
) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	hello := h.service.GetHello()

	rh.JSONResponse(
		hello,
		http.StatusOK,
	)
}
