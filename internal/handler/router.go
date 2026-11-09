package handler

import (
	"currency-converter/internal/handler/middleware"
	v1 "currency-converter/internal/handler/v1"
	"net/http"
)

type Router struct {
	authHandler *v1.Handler
}

func NewRouter(authHandler *v1.Handler) *Router {
	return &Router{authHandler: authHandler}
}

func (r *Router) InitRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/register", r.authHandler.Register)
	mux.HandleFunc("POST /api/v1/login", r.authHandler.Login)
	return middleware.Recover(mux)
}
