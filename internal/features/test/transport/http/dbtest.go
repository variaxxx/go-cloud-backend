package test_http

import (
	core_logger "cloud/internal/core/logger"
	core_http_request "cloud/internal/core/transport/http/request"
	core_http_response "cloud/internal/core/transport/http/response"
	"net/http"
)

type DbTestRequest struct {
	Str string `json:"str" validate:"required"`
}

func (h *Handler) DbTest(rw http.ResponseWriter, r *http.Request) {
	log := core_logger.FromContext(r.Context())
	rh := core_http_response.NewHTTPResponseHandler(log, rw)

	var request DbTestRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		rh.ErrorResponse(err, "Invalid request body")
		return
	}

	if err := h.service.WriteStr(r.Context(), request.Str); err != nil {
		rh.ErrorResponse(err, "Test failed")
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
