package main

import (
	"log"
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sundayyogurt/product_service/config"
	"github.com/sundayyogurt/product_service/infra/db"
	"github.com/sundayyogurt/product_service/infra/queue"
	"github.com/sundayyogurt/product_service/internal/api/grpc/product"
	"github.com/sundayyogurt/product_service/internal/api/rest"
	"github.com/sundayyogurt/product_service/internal/domain"
	"github.com/sundayyogurt/product_service/internal/repository"
	"github.com/sundayyogurt/product_service/internal/services"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := db.Connect(cfg.DatabaseDSN)

	// Run Migrations
	err = db.AutoMigrate(
		&domain.Category{},
		&domain.Product{},
	) // we will add out models here later

	// Initialize repository & services
	catalogRepo := repository.NewCatalogRepository(db)
	catalogService := services.NewCatalogService(catalogRepo)

	//start Kafka Consumer
	kafkaConsumer := queue.NewKafkaConsumer(cfg.KafkaBroker, "product-update-topic", "catalog-service-group", catalogService)
	go kafkaConsumer.Listen()

	// setup our GRPC
	go func() {
		list, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("failed to listen to GRPC Server: %v", err)
		}
		grpcServer := grpc.NewServer()
		productServer := product.NewProductGRPCService(catalogService)
		product.RegisterProductServiceServer(grpcServer, productServer)

		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve GRPC Server: %v", err)
		}

		log.Println("GRPC Server started on port 50052")
		// add some GRPC server code so we can respond to GRPC requests
	}()

	// setup Fiber app
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // allow all origins
		AllowHeaders: "Content-Type, Accept, Authorization",
	}))

	handler := rest.NewCatalogHandler(catalogService)
	handler.SetupRoutes(app)

	app.Get("/", HealthCheck)

	app.Listen("localhost:9000")
}

func HealthCheck(ctx *fiber.Ctx) error {
	return ctx.Status(200).JSON(fiber.Map{
		"message": "Healthy",
	})
}
