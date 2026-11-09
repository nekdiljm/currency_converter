package v1

import (
	"currency-converter/internal/handler/response"
	"currency-converter/internal/service"
	"currency-converter/pkg/serializer"
	"net/http"
)

type Handler struct {
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	req := registerRequest{}
	if err := serializer.Decode(r.Body, &req); err != nil {
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
		response.Error(w, err)
		return
	}

	tokens, err := h.authService.Login(r.Context(), req.toInput())
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, toLoginResponse(tokens), http.StatusOK)
}
