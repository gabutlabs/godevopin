package http_handler

import (
	"fmt"

	service "github.com/gabutlabs/godevopin/internal/services"
	"github.com/gabutlabs/godevopin/pkg"
	"github.com/gofiber/fiber/v2"
)

type SystemMetricHandler struct {
	systemMetricService service.SystemMetricService
}

func NewSystemMetricHandler(svc service.SystemMetricService) *SystemMetricHandler {
	return &SystemMetricHandler{systemMetricService: svc}
}

// GetSystemMetrics menghandle request GET /system-metrics
func (h *SystemMetricHandler) GetSystemMetrics(c *fiber.Ctx) error {
	metrics, err := h.systemMetricService.GetSystemMetrics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve system metrics"})
	}
	return c.JSON(pkg.GenerateResponse("System metrics retrieved successfully", metrics))
}

func (h *SystemMetricHandler) GetFilteredSystemMetrics(c *fiber.Ctx) error {
	filter := c.Query("filter", "1h") // Default filter 1 hour
	metrics, err := h.systemMetricService.GetFilteredSystemMetrics(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve filtered system metrics"})
	}
	return c.JSON(pkg.GenerateResponse("Filtered system metrics retrieved successfully", metrics))
}

func (h *SystemMetricHandler) GetCurrentDiskUsage(c *fiber.Ctx) error {
	diskUsage, err := h.systemMetricService.GetDiskUsage()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to get current disk usage"})
	}
	return c.JSON(pkg.GenerateResponse("Current disk usage successfully", diskUsage))
}

// SetupSystemMetricRoutes mendaftarkan rute metrik sistem ke router Fiber
func (h *SystemMetricHandler) SetupSystemMetricRoutes(router fiber.Router) {
	metricsGroup := router.Group("/system-metrics")
	metricsGroup.Get("/", h.GetSystemMetrics)
	metricsGroup.Get("/filtered", h.GetFilteredSystemMetrics)
	metricsGroup.Get("/disk-usage", h.GetCurrentDiskUsage)
}
