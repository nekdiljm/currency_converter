package v1

import (
	"currency-converter/internal/api/response"
	"currency-converter/pkg/serializer"
	"net/http"

	"go.uber.org/zap"
)

func (h *Handler) Converter(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("currency converter request", zap.Any("request", r.Body))
	req := convertRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	convertResp, err := h.currencyConverter.Convert(r.Context(), req.toInput())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toConvertResponse(convertResp), http.StatusOK)
}
