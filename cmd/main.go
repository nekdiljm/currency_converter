package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"currency-converter/internal/config"
	"currency-converter/internal/handler"
	"currency-converter/internal/handler/v1"
	"currency-converter/internal/repository/session"
	"currency-converter/internal/repository/user"
	"currency-converter/internal/service/auth"
	zlog "currency-converter/pkg/logger"
	"currency-converter/pkg/postgres"
	rdb "currency-converter/pkg/redis"

	"go.uber.org/zap"
)

const gracefulWait = 500 * time.Millisecond

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to init config %v", err)
	}

	logger, err := zlog.New(cfg.LoggerConfig.LogLevel)
	if err != nil {
		zap.L().Fatal("failed to init log", zap.Error(err))
	}

	ctx := context.Background()

	logger.Info("initialize db")
	db, err := postgres.New(ctx, cfg.PostgresConfig.ConnDb)
	if err != nil {
		logger.Fatal("failed to connect to postgres", zap.Error(err))
	}

	logger.Info("initialize redis")
	redisClient, err := rdb.New(ctx, cfg.RedisConfig.RedisHost, cfg.RedisConfig.RedisPort)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	userRepo := user.NewRepository(db)
	sessionRepo := session.NewRepository(redisClient)

	authService := auth.NewService(&cfg, userRepo, sessionRepo, logger)
	authHandler := v1.NewHandler(authService, logger)

	server := http.Server{
		Addr:    net.JoinHostPort(cfg.ServerConfig.HTTPHost, cfg.ServerConfig.HTTPPort),
		Handler: handler.NewRouter(authHandler, authService, logger).InitRouter(),
	}

	done := make(chan struct{}, 1)

	go func() {
		logger.Info("initialize server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", zap.Error(err))
			done <- struct{}{}
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signals:
	case <-done:
	}

	logger.Info("shutdown server")
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", zap.Error(err))
	}

	logger.Info("close redis")
	if err := redisClient.Close(); err != nil {
		logger.Error("failed to close redis", zap.Error(err))
	}

	logger.Info("close db")
	db.Close()

	logger.Info("sync logger")
	if err := logger.Sync(); err != nil {
		logger.Error("failed to sync logger", zap.Error(err))
	}

	time.Sleep(gracefulWait)
}
