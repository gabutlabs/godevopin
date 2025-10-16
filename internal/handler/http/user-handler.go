package http_handler

import (
	"strconv"
	"time"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
)

// UserHandler sekarang menggunakan interface UserService dari package service
type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{userService: svc}
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve users", nil))
	}
	return c.JSON(pkg.GenerateResponse("Users retrieved successfully", users))
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid user ID", nil))
	}
	user, err := h.userService.GetUserProfile(uint(id))
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("User not found", nil))
	}
	return c.JSON(pkg.GenerateResponse("User retrieved successfully", user))
}

// CreateUser menghandle request POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// DTO (Data Transfer Object) untuk request body
	type CreateUserRequest struct {
		Name     string `json:"name" validate:"required,min=2"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password"`
	}
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if errors := pkg.ValidateAndFormat(&req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Validation failed", errors))
	}

	if req.Password == "" {
		req.Password = "Password123!"
	}

	user := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.userService.RegisterUser(user.Name, user.Email, user.Password); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to create user", nil))
	}
	// Mengembalikan user yang baru dibuat (tanpa password)
	return c.Status(fiber.StatusCreated).JSON(pkg.GenerateResponse("User created successfully", map[string]any{"name": user.Name, "email": user.Email, "created_at": time.Now(), "updated_at": time.Now()}))
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	type UpdateUserRequest struct {
		Name  string `json:"name" validate:"required,min=2"`
		Email string `json:"email" validate:"required,email"`
	}
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid param ID", nil))
	}
	if errors := pkg.ValidateAndFormat(&req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Validation failed", errors))
	}
	currentUser, err := h.userService.GetUserProfile(uint(id))
	if err != nil || currentUser == nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("User not found", nil))
	}
	user := &model.User{
		ID:       uint(id),
		Name:     req.Name,
		Email:    req.Email,
		Password: currentUser.Password, // Pertahankan password lama
	}
	if err := h.userService.UpdateUserProfile(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update user", nil))
	}
	return c.JSON(pkg.GenerateResponse("User updated successfully", user))
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid param ID", nil))
	}
	if err := h.userService.DeleteUser(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to delete user", nil))
	}
	return c.JSON(pkg.GenerateResponse("User deleted successfully", nil))
}

// SetupUserRoutes mendaftarkan rute ke router Fiber
func (h *UserHandler) SetupUserRoutes(router fiber.Router) {
	userGroup := router.Group("/users")
	userGroup.Post("/", h.CreateUser)
	userGroup.Get("/", h.GetAllUsers)
	userGroup.Get("/:id", h.GetUserByID)
	userGroup.Put("/:id", h.UpdateUser)
	userGroup.Delete("/:id", h.DeleteUser)
}
