package api

import (
	"currency-converter/internal/service"
	"net/http"

	"currency-converter/internal/api/middleware"
	v1 "currency-converter/internal/api/v1"
	"currency-converter/pkg/logger"
)

type Router struct {
	handler     *v1.Handler
	authService service.AuthService
	logger      logger.Logger
}

func NewRouter(handler *v1.Handler, authService service.AuthService, logger logger.Logger) *Router {
	return &Router{
		handler:     handler,
		authService: authService,
		logger:      logger,
	}
}

func (r *Router) InitRouter() http.Handler {
	mux := http.NewServeMux()
	mw := middleware.New(r.authService, r.logger)

	mux.HandleFunc("POST /api/v1/register", r.handler.Register)
	mux.HandleFunc("POST /api/v1/login", r.handler.Login)
	mux.HandleFunc("POST /api/v1/logout", r.handler.Logout)
	mux.HandleFunc("POST /api/v1/refresh", r.handler.Refresh)

	mux.Handle("POST /api/v1/currency-converter", mw.ApplyAuth(r.handler.Convert))
	mux.Handle("GET /api/v1/conversion-history", mw.ApplyAuth(r.handler.GetHistoryByUserID))
	mux.Handle("GET /api/v1/currency-rates", mw.ApplyAuth(r.handler.GetCurrencyRatesByDate))

	mux.Handle("POST /api/v1/create-favorite-pair", mw.ApplyAuth(r.handler.CreateFavoritePair))
	mux.Handle("GET /api/v1/favorite-pairs", mw.ApplyAuth(r.handler.GetFavoritePairs))
	mux.Handle("POST /api/v1/delete-favorite-pair", mw.ApplyAuth(r.handler.DeleteFavoritePair))

	return mw.Logging(middleware.Recover(mux))
}
