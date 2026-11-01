package main

import (
	"context"
	"currency-converter/internal/config"
	"currency-converter/pkg/postgres"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Println("error to load config", err)
		return
	}

	ctx := context.Background()

	dbConn, err := postgres.New(ctx, cfg.ConnDb)
	if err != nil {
		log.Println("error connect to db", err)
		return
	}
	defer dbConn.Close()

}
