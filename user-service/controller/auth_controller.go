package controller

import (
	"micro-warehouse/user-service/controller/request"
	"micro-warehouse/user-service/controller/response"
	"micro-warehouse/user-service/pkg/conv"
	"micro-warehouse/user-service/pkg/validator"
	"micro-warehouse/user-service/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AuthControllerInterface interface {
	Login(c *fiber.Ctx) error
}

type AuthController struct {
	AuthService usecase.UserUsecaseInterface
}

func NewAuthController(authService usecase.UserUsecaseInterface) AuthControllerInterface {
	return &AuthController{
		AuthService: authService,
	}
}

// Login implements [AuthControllerInterface].
func (a *AuthController) Login(c *fiber.Ctx) error {
	ctx := c.Context()

	req := request.LoginRequest{}
	if err := c.BodyParser(&req); err != nil {
		log.Errorf("[AuthController] Login - 1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validator.Validate(req); err != nil {
		log.Errorf("[AuthController] Login - 2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := a.AuthService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Errorf("[AuthController] Login - 3: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if user == nil {
		log.Errorf("[AuthController] Login - 4: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	isSame := conv.CheckPasswordHash(req.Password, user.Password)
	if !isSame {
		log.Errorf("[AuthController] Login - 5: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	var roles []string

	for _, r := range user.Roles {
		roles = append(roles, r.Name)
	}

	resp := response.LoginResponse{
		UserID: user.ID,
		Email:  user.Email,
		Role:   roles,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    resp,
		"message": "Login successful",
	})
}
