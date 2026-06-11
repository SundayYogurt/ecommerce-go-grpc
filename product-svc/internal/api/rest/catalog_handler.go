package rest

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sundayyogurt/product_service/internal/dto"
	"github.com/sundayyogurt/product_service/internal/services"
)

type CatalogHandler struct {
	svc services.CatalogService
}

func NewCatalogHandler(svc services.CatalogService) *CatalogHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) SetupRoutes(app *fiber.App) {
	app.Post("/categories", h.CreateCategory)
	app.Get("/categories", h.GetAllCategories)
	app.Get("/categories/:id", h.GetCategoryByID)
	app.Patch("/categories/:id", h.UpdateCategory)

	app.Post("/products", h.CreateProduct)
	app.Get("/products", h.GetAllProducts)
	app.Get("/products/:id", h.GetProductByID)
	app.Patch("/products/:id", h.UpdateProduct)
}

func (h *CatalogHandler) CreateCategory(ctx *fiber.Ctx) error {
	var input dto.CreateCategoryRequest
	if err := ctx.BodyParser(&input); err != nil {
		log.Printf("createCategory failed to parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.svc.CreateCategory(input); err != nil {
		log.Printf("Error creating category: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Category created successfully"})
}

func (h *CatalogHandler) GetAllCategories(ctx *fiber.Ctx) error {
	categories, err := h.svc.GetAllCategories()
	if err != nil {
		log.Printf("Error getting all categories: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(categories)
}

func (h *CatalogHandler) GetCategoryByID(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	category, err := h.svc.GetCategoryByID(uint(id))
	if err != nil || category == nil {
		log.Printf("Error getting category by ID: %v", err)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(category)
}

func (h *CatalogHandler) UpdateCategory(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	var input dto.UpdateCategoryRequest

	if err := ctx.BodyParser(&input); err != nil {
		log.Printf("updateCategory failed to parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.svc.UpdateCategory(uint(id), input); err != nil {
		log.Printf("Error updating category: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Category updated successfully"})
}

func (h *CatalogHandler) CreateProduct(ctx *fiber.Ctx) error {

	var input dto.CreateProductRequest

	if err := ctx.BodyParser(&input); err != nil {
		log.Printf("createProduct failed to parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.svc.CreateProduct(input); err != nil {
		log.Printf("Error creating product: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Product created successfully"})
}

func (h *CatalogHandler) GetAllProducts(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	products, err := h.svc.GetProductByID(uint(id))
	if err != nil {
		log.Printf("Error getting product by ID: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(products)
}

func (h *CatalogHandler) GetProductByID(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	product, err := h.svc.GetProductByID(uint(id))
	if err != nil {
		log.Printf("Error getting product by ID: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if product == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	return ctx.Status(fiber.StatusOK).JSON(product)
}

func (h *CatalogHandler) UpdateProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	var input dto.UpdateProductRequest

	if err := ctx.BodyParser(&input); err != nil {
		log.Printf("updateProduct failed to parsing request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.svc.UpdateProduct(uint(id), input); err != nil {
		log.Printf("Error updating product: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product updated successfully"})
}
