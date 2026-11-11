package api

import (
	"currency-converter/internal/service/auth"
	"net/http"

	"currency-converter/internal/api/middleware"
	v1 "currency-converter/internal/api/v1"
	"currency-converter/pkg/logger"
)

type Router struct {
	handler     *v1.Handler
	authService auth.AuthService
	logger      logger.Logger
}

func NewRouter(handler *v1.Handler, authService auth.AuthService, logger logger.Logger) *Router {
	return &Router{
		handler:     handler,
		authService: authService,
		logger:      logger,
	}
}

func (r *Router) InitRouter() http.Handler {
	mux := http.NewServeMux()
	mw := middleware.New(r.authService, r.logger)
	mwFn := mw.Auth()

	mux.HandleFunc("POST /api/v1/register", r.handler.Register)
	mux.HandleFunc("POST /api/v1/login", r.handler.Login)
	mux.HandleFunc("POST /api/v1/logout", r.handler.Logout)
	mux.HandleFunc("POST /api/v1/refresh", r.handler.Refresh)

	mux.Handle("QUERY /api/v1/currency-converter", mwFn(http.HandlerFunc(r.handler.Converter)))

	return middleware.Recover(mux)
}
