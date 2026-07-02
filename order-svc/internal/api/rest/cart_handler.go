package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sundayyogurt/order_service/internal/dto"
	"github.com/sundayyogurt/order_service/internal/services"
)

type CartHandler struct {
	svc services.CartService
}

func NewCartHandler(s services.CartService) *CartHandler {
	return &CartHandler{svc: s}
}

func (h *CartHandler) Register(r fiber.Router) {
	r.Post("/cart", h.add)
	r.Get("/cart", h.get)
	r.Patch("/cart", h.updateQty)
	r.Delete("/cart/:productId", h.remove)
}

func (h *CartHandler) fetchAuthorizedUserId(c *fiber.Ctx) (uint, error) {
	rawUserData := c.Get("X-User-Id")
	userId, err := strconv.Atoi(rawUserData)
	if err != nil || userId <= 0 {
		return 0, fiber.NewError(fiber.StatusBadRequest, "missing or invalid  X-User-Id header")
	}
	return uint(userId), nil
}

func (h *CartHandler) add(c *fiber.Ctx) error {
	uid, err := h.fetchAuthorizedUserId(c)
	if err != nil {
		return err
	}

	var CartReq dto.CartAddRequest
	if err := c.BodyParser(&CartReq); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.svc.Add(uid, CartReq); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *CartHandler) get(c *fiber.Ctx) error {
	uid, err := h.fetchAuthorizedUserId(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.Get(uid)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(resp)
}

func (h *CartHandler) updateQty(c *fiber.Ctx) error {
	uid, err := h.fetchAuthorizedUserId(c)
	if err != nil {
		return err
	}

	var req dto.CartAddRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid input")
	}

	if err := h.svc.UpdateQty(uid, req.ProductID, req.Qty); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"status": "ok"})
}

func (h *CartHandler) remove(c *fiber.Ctx) error {
	uid, err := h.fetchAuthorizedUserId(c)
	if err != nil {
		return err
	}

	var productInput = c.Params("productId")
	pid, _ := strconv.Atoi(productInput)

	if pid <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid productId parameter")
	}

	if err := h.svc.Remove(uid, uint(pid)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{"status": "ok"})
}
