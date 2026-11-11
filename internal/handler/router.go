package handler

import (
	"currency-converter/internal/service/auth"
	"net/http"

	"currency-converter/internal/handler/middleware"
	v1 "currency-converter/internal/handler/v1"
	"currency-converter/pkg/logger"
)

type Router struct {
	authHandler *v1.Handler
	authService auth.AuthService
	logger      logger.Logger
}

func NewRouter(authHandler *v1.Handler, authService auth.AuthService, logger logger.Logger) *Router {
	return &Router{
		authHandler: authHandler,
		authService: authService,
		logger:      logger,
	}
}

func (r *Router) InitRouter() http.Handler {
	mux := http.NewServeMux()
	//mw := middleware.New(r.authService, r.logger)
	//mwFn := mw.Auth()

	mux.HandleFunc("POST /api/v1/register", r.authHandler.Register)
	mux.HandleFunc("POST /api/v1/login", r.authHandler.Login)
	mux.HandleFunc("POST /api/v1/logout", r.authHandler.Logout)
	mux.HandleFunc("POST /api/v1/refresh", r.authHandler.Refresh)

	//mux.HandleFunc("POST /api/v1/currency-converter", mwFn(http.HandleFunc(r.authHandler.Refresh)))

	return middleware.Recover(mux)
}
