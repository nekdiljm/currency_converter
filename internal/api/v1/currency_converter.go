package v1

import (
	"currency-converter/internal/api/middleware"
	"currency-converter/internal/api/response"
	"currency-converter/internal/models"
	"currency-converter/pkg/serializer"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (h *Handler) Convert(w http.ResponseWriter, r *http.Request) {
	req := convertRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	h.logger.Debug("currency converter request", zap.Any("request", req))

	ctx := r.Context()
	userID := middleware.UserIDFromCtx(ctx)
	input := req.toInput()
	input.UserID = userID

	convertResp, err := h.services.CurrencyConverter.Convert(ctx, input)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toConvertResponse(convertResp), http.StatusOK)
}

func (h *Handler) GetHistoryByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromCtx(ctx)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	cursor, _ := strconv.Atoi(r.URL.Query().Get("cursor"))

	conversionHistory, err := h.services.CurrencyConverter.GetHistoryByUserID(ctx, userID, limit, cursor)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toHistoryResponse(conversionHistory), http.StatusOK)
}

func (h *Handler) GetCurrencyRatesByDate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fromDate := r.URL.Query().Get("from")
	toDate := r.URL.Query().Get("to")
	baseCurrency := r.URL.Query().Get("base")

	currencyRatesByDateResp, err := h.services.CurrencyConverter.GetCurrencyRatesByDate(ctx, models.CurrencyRatesByDateRequest{
		FromDate:     fromDate,
		ToDate:       toDate,
		BaseCurrency: baseCurrency,
	})
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toCurrencyRatesByDateResponse(currencyRatesByDateResp), http.StatusOK)
}
