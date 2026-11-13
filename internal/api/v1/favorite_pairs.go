package v1

import (
	"currency-converter/internal/api/middleware"
	"currency-converter/internal/api/response"
	"currency-converter/pkg/serializer"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

func (h *Handler) CreateFavoritePair(w http.ResponseWriter, r *http.Request) {
	req := favoritePairCreateRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	ctx := r.Context()
	userID := middleware.UserIDFromCtx(ctx)
	favPair := req.toInput()
	favPair.UserID = userID

	id, err := h.favoritePairs.CreateFavPair(ctx, favPair)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toFavoritePairCreateResponse(id), http.StatusCreated)
}

func (h *Handler) GetFavoritePairs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromCtx(ctx)
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	favoritePairs, err := h.favoritePairs.GetFavoritePairByUserID(ctx, userID, limit, offset)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toFavoritePairsResponse(favoritePairs), http.StatusOK)
}

func (h *Handler) DeleteFavoritePair(w http.ResponseWriter, r *http.Request) {
	req := favoritePairsDeleteRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	ctx := r.Context()
	if err := h.favoritePairs.DeleteFavoritePair(ctx, req.ID); err != nil {
		response.Error(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
