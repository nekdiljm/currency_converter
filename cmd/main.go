package main

import (
	"context"
	"currency-converter/internal/client/currencyratesprovider"
	"currency-converter/internal/service/currencyconverter"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"currency-converter/internal/api"
	"currency-converter/internal/api/v1"
	"currency-converter/internal/config"
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

	currencyRatesProvider := currencyratesprovider.NewClient(cfg.CurrencyRatesConfig.Host, cfg.CurrencyRatesConfig.Timeout)
	currencyConverterService := currencyconverter.NewService(currencyRatesProvider, redisClient, logger)

	authService := auth.NewService(&cfg, userRepo, sessionRepo, logger)
	handler := v1.NewHandler(authService, currencyConverterService, logger)

	server := http.Server{
		Addr:    net.JoinHostPort(cfg.ServerConfig.HTTPHost, cfg.ServerConfig.HTTPPort),
		Handler: api.NewRouter(handler, authService, logger).InitRouter(),
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

	logger.Info("close currency rates provider")
	currencyRatesProvider.Close()

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
