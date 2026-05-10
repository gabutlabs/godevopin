package http_handler

import (
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gofiber/fiber/v2"
)

type SettingHandler struct {
	service service.SettingService
}

func NewSettingHandler(service service.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (h *SettingHandler) SetupSettingRoutes(router fiber.Router) {
	group := router.Group("/settings")
	group.Get("/", h.GetSettings)
	group.Put("/", h.UpdateSettings)
}

func (h *SettingHandler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetSettings()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(settings)
}

func (h *SettingHandler) UpdateSettings(c *fiber.Ctx) error {
	var req service.UpdateSettingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if err := h.service.UpdateSettings(&req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Settings updated successfully",
	})
}
