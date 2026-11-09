package main

import (
	"context"
	"currency-converter/internal/config"
	"currency-converter/internal/handler"
	v1 "currency-converter/internal/handler/v1"
	"currency-converter/internal/repository/user"
	"currency-converter/internal/service/auth"
	"currency-converter/pkg/postgres"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const gracefulWait = 500 * time.Millisecond

func main() {
	// TODO implement logger -zap.Logger | slog | logrus and etc
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config %v", err)
	}

	ctx := context.Background()

	db, err := postgres.New(ctx, cfg.ConnDb)
	if err != nil {
		log.Fatalf("connect to db %v", err)
	}

	userRepo := user.NewRepository(db)
	authService := auth.NewService(userRepo)
	authHandler := v1.NewHandler(authService)

	server := http.Server{
		Addr:    cfg.HTTPHost + ":" + cfg.HTTPPort,
		Handler: handler.NewRouter(authHandler).InitRouter(),
	}

	fmt.Printf("App started on: %s", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			return
		}
	}()

	// TODO
	// redis init
	// main blocks forever

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	// TODO graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		// логировать через zap
	}
	db.Close()

	fmt.Println("server shutdown")

	time.Sleep(gracefulWait)
}
