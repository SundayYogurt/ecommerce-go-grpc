package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sundayyogurt/user_service/config"
	"github.com/sundayyogurt/user_service/infra/queue"
	"github.com/sundayyogurt/user_service/internal/api/grpc/order"
)

func main() {

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // allow all origins
		AllowHeaders: "Content-Type, Accept, Authorization",
	}))

	// Load config
	cfg := config.LoadConfig()

	kafkaProducer := queue.NewProducer(cfg.KafkaBroker, cfg.KafkaTopic)
	log.Printf("Kafka Producer created %v", kafkaProducer)

	orderServiceClient, err := order.NewOrderServiceGRPCClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to conneect to order service: &v", err)
	}

	log.Printf("Order Service created %v", orderServiceClient)

	app.Get("/", HealthCheck)

	app.Listen("localhost:9000")
}

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "Healthy",
	})
}
