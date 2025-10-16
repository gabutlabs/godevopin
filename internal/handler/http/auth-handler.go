package http_handler

import (
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(svc service.UserService) *AuthHandler {
	return &AuthHandler{userService: svc}
}

func (h *AuthHandler) LoginHandler(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	data, err := h.userService.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(pkg.GenerateResponse("Login successful", data))
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	getUser, err := h.userService.GetUserProfile(uint(claims["id"].(float64)))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(pkg.GenerateResponse("User info retrieved successfully", fiber.Map{"is_loggedin": true, "user": getUser}))
}

// SetupAuthRoutes mendaftarkan rute autentikasi ke router Fiber
func (h *AuthHandler) SetupAuthRoutes(router fiber.Router) {
	authGroup := router.Group("/auth")
	authGroup.Post("/login", h.LoginHandler)
	// Tambahkan rute lain seperti register, logout, dll jika diperlukan
}
