package v1

import (
	"currency-converter/internal/service/auth"
	"net/http"

	"go.uber.org/zap"

	"currency-converter/internal/handler/response"
	"currency-converter/pkg/logger"
	"currency-converter/pkg/serializer"
)

type Handler struct {
	authService auth.AuthService
	logger      logger.Logger
}

func NewHandler(authService auth.AuthService, logger logger.Logger) *Handler {
	return &Handler{
		authService: authService,
		logger:      logger,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := registerRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		response.Error(w, err)
		return
	}

	user, err := h.authService.Register(r.Context(), req.toInput())
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

	tokens, err := h.authService.Login(r.Context(), req.toInput())
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

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
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

	tokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toTokenPairsResponse(tokens), http.StatusOK)
}
