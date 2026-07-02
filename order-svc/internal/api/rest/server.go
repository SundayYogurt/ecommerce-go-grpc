package rest

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sundayyogurt/order_service/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	app         *fiber.App
	port        string
	db          *gorm.DB
	cartService services.CartService
}

func NewServer(port string, cartSvc services.CartService, db *gorm.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName: "order-service",
	})
	return &Server{
		app:         app,
		port:        port,
		cartService: cartSvc,
		db:          db,
	}
}

func (s *Server) registerHandlers() {
	// Liveness: check the container is running correctly
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// readiness: check if service is ready to accept traffic
	s.app.Get("/ready", func(c *fiber.Ctx) error {
		db, err := s.db.DB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("DB Error")
		}
		if err := db.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("db not ready")
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api := s.app.Group("/api/v1")

	cartHandler := NewCartHandler(s.cartService)
	cartHandler.Register(api)

}

func (s *Server) Start() error {
	s.registerHandlers()
	log.Printf("Rest APIs are Listening on port %s", s.port)
	return s.app.Listen(s.port)
}
