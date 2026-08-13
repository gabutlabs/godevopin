package http_handler

import (
	"strconv"

	service "github.com/gabutlabs/godevopin/internal/services"
	"github.com/gabutlabs/godevopin/pkg"
	"github.com/gofiber/fiber/v2"
)

type ProcessMonitoringHandler struct {
	service service.ProcessMonitoringService
}

func NewProcessMonitoringHandler(svc service.ProcessMonitoringService) *ProcessMonitoringHandler {
	return &ProcessMonitoringHandler{service: svc}
}

func (h *ProcessMonitoringHandler) GetLiveProcesses(c *fiber.Ctx) error {
	limit := 200
	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid process limit", nil))
		}
		limit = parsed
	}

	processes, err := h.service.GetLiveProcesses(c.Query("sort", "cpu"), c.Query("search"), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve live processes", nil))
	}
	return c.JSON(pkg.GenerateResponse("Live processes retrieved successfully", processes))
}

func (h *ProcessMonitoringHandler) GetProcessHistory(c *fiber.Ctx) error {
	processKey := c.Query("process_key")
	if processKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("process_key is required", nil))
	}

	history, err := h.service.GetProcessHistory(processKey, c.Query("filter", "1h"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("Process history retrieved successfully", history))
}

func (h *ProcessMonitoringHandler) SetupProcessMonitoringRoutes(router fiber.Router) {
	processesGroup := router.Group("/processes")
	processesGroup.Get("/live", h.GetLiveProcesses)
	processesGroup.Get("/history", h.GetProcessHistory)
}
