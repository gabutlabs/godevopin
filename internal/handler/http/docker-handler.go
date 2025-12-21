package http_handler

import (
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
)

type DockerHandler struct {
	service *service.DockerService
}

// NewDockerHandler creates a new Docker handler
func NewDockerHandler(dockerService *service.DockerService) *DockerHandler {
	return &DockerHandler{
		service: dockerService,
	}
}

// GetContainers handles GET /api/containers
func (h *DockerHandler) GetContainers(c *fiber.Ctx) error {
	all := c.QueryBool("all", true)

	containers, err := h.service.GetContainers(c.Context(), all)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			pkg.GenerateErrorResponse("Failed to get containers", err.Error()),
		)
	}

	return c.JSON(pkg.GenerateResponse("Containers retrieved successfully", containers))
}

// GetImages handles GET /api/images
func (h *DockerHandler) GetImages(c *fiber.Ctx) error {
	images, err := h.service.GetImages(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			pkg.GenerateErrorResponse("Failed to get images", err.Error()),
		)
	}

	return c.JSON(pkg.GenerateResponse("Images retrieved successfully", images))
}

// GetNetworks handles GET /api/networks
func (h *DockerHandler) GetNetworks(c *fiber.Ctx) error {
	networks, err := h.service.GetNetworks(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			pkg.GenerateErrorResponse("Failed to get networks", err.Error()),
		)
	}

	return c.JSON(pkg.GenerateResponse("Networks retrieved successfully", networks))
}

// GetVolumes handles GET /api/volumes
func (h *DockerHandler) GetVolumes(c *fiber.Ctx) error {
	volumes, err := h.service.GetVolumes(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			pkg.GenerateErrorResponse("Failed to get volumes", err.Error()),
		)
	}

	return c.JSON(pkg.GenerateResponse("Volumes retrieved successfully", volumes))
}

// GetAllResources handles GET /api/resources
func (h *DockerHandler) GetAllResources(c *fiber.Ctx) error {
	resources, err := h.service.GetAllResources(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			pkg.GenerateErrorResponse("Failed to get resources", err.Error()),
		)
	}

	return c.JSON(pkg.GenerateResponse("Resources retrieved successfully", resources))
}

// SetupAlarmRoutes registers the alarm routes to the provided router
func (dh *DockerHandler) SetupDockerRoutes(router fiber.Router) {
	docker := router.Group("/docker")

	docker.Get("/containers", dh.GetContainers)
	docker.Get("/images", dh.GetImages)
	docker.Get("/networks", dh.GetNetworks)
	docker.Get("/volumes", dh.GetVolumes)
	docker.Get("/resources", dh.GetAllResources)

}
