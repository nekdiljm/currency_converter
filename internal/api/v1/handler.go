package v1

import (
	"currency-converter/internal/service"
	"currency-converter/pkg/logger"
)

type Handler struct {
	services service.Service
	logger   logger.Logger
}

func NewHandler(
	services service.Service,
	logger logger.Logger,
) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}
