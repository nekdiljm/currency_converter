package main

import (
	"context"
	"currency-converter/internal/config"
	"currency-converter/pkg/postgres"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config %v", err)
	}

	ctx := context.Background()

	dbConn, err := postgres.New(ctx, cfg.ConnDb)
	if err != nil {
		log.Fatalf("connect to db %v", err)
	}
	defer dbConn.Close()

	fmt.Println("App started")

}
