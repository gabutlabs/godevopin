package http_handler

import (
	"strconv"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
)

// WorkerServiceHandler handles WorkerService-related HTTP requests
type WorkerServiceHandler struct {
	workerServiceService service.WorkerServiceService
}

func NewWorkerServiceHandler(svc service.WorkerServiceService) *WorkerServiceHandler {
	return &WorkerServiceHandler{workerServiceService: svc}
}

// GetAllWorkerServices handles GET /worker-services request
func (h *WorkerServiceHandler) GetAllWorkerServices(c *fiber.Ctx) error {
	workerServices, err := h.workerServiceService.ListWorkerServices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve worker services", nil))
	}
	return c.JSON(pkg.GenerateResponse("Worker services retrieved successfully", workerServices))
}

// GetWorkerServiceByID handles GET /worker-services/:id request
func (h *WorkerServiceHandler) GetWorkerServiceByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid worker service ID", nil))
	}
	workerService, err := h.workerServiceService.GetWorkerServiceByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("Worker service not found", nil))
	}
	return c.JSON(pkg.GenerateResponse("Worker service retrieved successfully", workerService))
}

// GetWorkerServiceByName handles GET /worker-services/name/:name request
func (h *WorkerServiceHandler) GetWorkerServiceByName(c *fiber.Ctx) error {
	name := c.Params("name")
	workerService, err := h.workerServiceService.GetWorkerServiceByName(name)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("Worker service not found", nil))
	}
	return c.JSON(pkg.GenerateResponse("Worker service retrieved successfully", workerService))
}

// CreateWorkerService handles POST /worker-services request
func (h *WorkerServiceHandler) CreateWorkerService(c *fiber.Ctx) error {
	type CreateWorkerServiceRequest struct {
		Name        string                `json:"name" validate:"required,min=2"`
		Description string                `json:"description"`
		DesiredState model.DesiredState   `json:"desired_state" validate:"required,oneof=enabled disabled"`
	}

	var req CreateWorkerServiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if errors := pkg.ValidateAndFormat(&req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Validation failed", errors))
	}

	if req.DesiredState == "" {
		req.DesiredState = model.StateEnabled
	}

	err := h.workerServiceService.CreateWorkerService(req.Name, req.Description, req.DesiredState)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to create worker service", nil))
	}

	return c.Status(fiber.StatusCreated).JSON(pkg.GenerateResponse("Worker service created successfully", map[string]any{
		"name":          req.Name,
		"description":   req.Description,
		"desired_state": req.DesiredState,
	}))
}

// UpdateWorkerService handles PUT /worker-services/:id request
func (h *WorkerServiceHandler) UpdateWorkerService(c *fiber.Ctx) error {
	type UpdateWorkerServiceRequest struct {
		Name        string               `json:"name" validate:"required,min=2"`
		Description string               `json:"description"`
		DesiredState model.DesiredState  `json:"desired_state" validate:"required,oneof=enabled disabled"`
	}

	var req UpdateWorkerServiceRequest
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

	currentWorkerService, err := h.workerServiceService.GetWorkerServiceByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("Worker service not found", nil))
	}

	updatedWorkerService := &model.WorkerService{
		ID:           uint(id),
		Name:         req.Name,
		Description:  req.Description,
		DesiredState: req.DesiredState,
		// Keep existing status
		CurrentStatus: currentWorkerService.CurrentStatus,
		HealthStatus:  currentWorkerService.HealthStatus,
		// Keep existing timestamps
		CreatedAt: currentWorkerService.CreatedAt,
		UpdatedAt: currentWorkerService.UpdatedAt,
	}

	if err := h.workerServiceService.UpdateWorkerService(updatedWorkerService); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update worker service", nil))
	}

	return c.JSON(pkg.GenerateResponse("Worker service updated successfully", updatedWorkerService))
}

// DeleteWorkerService handles DELETE /worker-services/:id request
func (h *WorkerServiceHandler) DeleteWorkerService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid param ID", nil))
	}

	if err := h.workerServiceService.DeleteWorkerService(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to delete worker service", nil))
	}

	return c.JSON(pkg.GenerateResponse("Worker service deleted successfully", nil))
}

// UpdateWorkerServiceStatus handles PUT /worker-services/:id/status request
func (h *WorkerServiceHandler) UpdateWorkerServiceStatus(c *fiber.Ctx) error {
	type UpdateStatusRequest struct {
		CurrentStatus model.CurrentStatus `json:"current_status" validate:"required,oneof=starting running stopped failed degraded restart"`
		HealthStatus  model.HealthStatus  `json:"health_status" validate:"required,oneof=healthy unhealthy unknown"`
	}

	var req UpdateStatusRequest
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

	if err := h.workerServiceService.UpdateWorkerServiceStatus(uint(id), req.CurrentStatus, req.HealthStatus); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update worker service status", nil))
	}

	return c.JSON(pkg.GenerateResponse("Worker service status updated successfully", nil))
}

// UpdateWorkerServiceHeartbeat handles PUT /worker-services/:id/heartbeat request
func (h *WorkerServiceHandler) UpdateWorkerServiceHeartbeat(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid param ID", nil))
	}

	if err := h.workerServiceService.UpdateWorkerServiceHeartbeat(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update worker service heartbeat", nil))
	}

	return c.JSON(pkg.GenerateResponse("Worker service heartbeat updated successfully", nil))
}

// SetupWorkerServiceRoutes registers routes to the Fiber router
func (h *WorkerServiceHandler) SetupWorkerServiceRoutes(router fiber.Router) {
	workerServiceGroup := router.Group("/worker-services")
	workerServiceGroup.Post("/", h.CreateWorkerService)
	workerServiceGroup.Get("/", h.GetAllWorkerServices)
	workerServiceGroup.Get("/:id", h.GetWorkerServiceByID)
	workerServiceGroup.Get("/name/:name", h.GetWorkerServiceByName)
	workerServiceGroup.Put("/:id", h.UpdateWorkerService)
	workerServiceGroup.Put("/:id/status", h.UpdateWorkerServiceStatus)
	workerServiceGroup.Put("/:id/heartbeat", h.UpdateWorkerServiceHeartbeat)
	workerServiceGroup.Delete("/:id", h.DeleteWorkerService)
}