package v1

import (
	"net/http"

	"go.uber.org/zap"

	"currency-converter/internal/api/response"
	"currency-converter/pkg/serializer"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := registerRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	user, err := h.services.AuthService.Register(r.Context(), req.toInput())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toUserResponse(user), http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	req := loginRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	tokens, err := h.services.AuthService.Login(r.Context(), req.toInput())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toTokenPairsResponse(tokens), http.StatusOK)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	req := logoutRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	err := h.services.AuthService.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	req := refreshRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	tokens, err := h.services.AuthService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toTokenPairsResponse(tokens), http.StatusOK)
}
