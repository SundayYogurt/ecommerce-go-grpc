package main

import (
	"log"

	"github.com/sundayyogurt/order_service/config"
	"github.com/sundayyogurt/order_service/infra/db/db"
	"github.com/sundayyogurt/order_service/internal/api/rest"
	"github.com/sundayyogurt/order_service/internal/repository"
	"github.com/sundayyogurt/order_service/internal/services"
)

func main() {

	cfg := config.LoadConfig()

	dbConn, err := db.Connect(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	cartRepository := repository.NewCartRepository(dbConn)
	cartService := services.NewCartService(cartRepository)

	server := rest.NewServer(cfg.ServerPort, cartService, dbConn)
	err = server.Start()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
