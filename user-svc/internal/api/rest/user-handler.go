package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sundayyogurt/user_service/internal/api/rest/middlewares"
	"github.com/sundayyogurt/user_service/internal/dto"
	"github.com/sundayyogurt/user_service/internal/services"
	"github.com/sundayyogurt/user_service/pkg/utils"
)

type UserHandler struct {
	svc services.UserService // แค่บอกว่า "มี field ชื่อ svc"
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SetupRoutes(app *fiber.App) {
	app.Post("/register", h.Register)
	app.Post("/login", h.Login)
	app.Post("/forgot-password", h.ForgotPassword)
	app.Post("/set-password", h.SetPassword)

	app.Use(middlewares.AuthMiddleware())
	app.Post("/profile", h.CreateProfile)
	app.Get("/profile", h.GetProfile)
	app.Get("/auth", h.Auth)
	app.Get("/me", h.Me)
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var reqBody dto.UserSignup

	if err := ctx.BodyParser(&reqBody); err != nil {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "please provide valid inputs")
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "email and password are required")
	}

	if err := h.svc.Register(reqBody); err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusCreated, "user registered successfully")
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var reqBody dto.UserLogin
	if err := ctx.BodyParser(&reqBody); err != nil {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "invalid inputs")
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "email and password are required")
	}

	user, err := h.svc.Login(reqBody.Email, reqBody.Password)
	if err != nil {
		return utils.RespondError(ctx, fiber.StatusUnauthorized, "invalid email or password")
	}

	token, err := middlewares.GenerateToken(int(user.ID), user.Email)

	if err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, "could not generate token")
	}
	return utils.RespondSuccess(ctx, fiber.StatusOK, fiber.Map{
		"token": token,
	})
}

func (h *UserHandler) ForgotPassword(ctx *fiber.Ctx) error {
	var reqBody struct {
		Email string `json:"email"`
	}

	if err := ctx.BodyParser(&reqBody); err != nil || reqBody.Email == "" {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "please provide a valid email id")
	}

	if err := h.svc.ForgotPassword(reqBody.Email); err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusOK, "password reset link sent to your email")
}

func (h *UserHandler) SetPassword(ctx *fiber.Ctx) error {
	var reqBody struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&reqBody); err != nil || reqBody.Token == "" || reqBody.Password == "" {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "please provide valid inputs")
	}

	if err := h.svc.SetPassword(reqBody.Token, reqBody.Password); err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusOK, "password reset link sent")
}

func (h *UserHandler) CreateProfile(ctx *fiber.Ctx) error {
	var reqBody dto.UserProfile
	if err := ctx.BodyParser(&reqBody); err != nil {
		return utils.RespondError(ctx, fiber.StatusBadRequest, "invalid inputs")
	}

	userID := ctx.Locals("user_id")
	reqBody.UserID = userID.(int)

	if err := h.svc.CreateProfile(reqBody); err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusOK, "profile created successfully")
}

func (h *UserHandler) GetProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("user_id").(int)
	profile, err := h.svc.GetProfile(userId)
	if err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.RespondSuccess(ctx, fiber.StatusOK, profile)
}

func (h *UserHandler) Auth(ctx *fiber.Ctx) error {
	user, err := h.svc.Authenticate(ctx)
	if err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusOK, fiber.Map{
		"authenticated": true,
		"user":          user,
	})
}

func (h *UserHandler) Me(ctx *fiber.Ctx) error {

	userID := ctx.Locals("user_id").(int)

	user, err := h.svc.GetProfile(userID)
	if err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	// fetch Orders via GRPC
	if err != nil {
		return utils.RespondError(ctx, fiber.StatusInternalServerError, err.Error())
	}

	return utils.RespondSuccess(ctx, fiber.StatusOK, fiber.Map{
		"user": user,
	})
}
