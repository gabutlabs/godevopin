package http_handler

import (
	"strconv"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectService service.ProjectService
	logService     service.LogHistoryService
}

func NewProjectHandler(svc service.ProjectService, ls service.LogHistoryService) *ProjectHandler {
	return &ProjectHandler{
		projectService: svc,
		logService:     ls,
	}
}

func (h *ProjectHandler) GetAllProjects(c *fiber.Ctx) error {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve projects", nil))
	}
	return c.JSON(pkg.GenerateResponse("Projects retrieved successfully", projects))
}

func (h *ProjectHandler) GetProjectByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid project ID", nil))
	}
	project, err := h.projectService.GetProjectByID(uint(id))
	if err != nil || project == nil {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("Project not found", nil))
	}
	return c.JSON(pkg.GenerateResponse("Project retrieved successfully", project))
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	type CreateProjectRequest struct {
		Name        string `json:"name" validate:"required,min=2"`
		PathLog     string `json:"path_log" validate:"required"`
		ProjectType string `json:"project_type" validate:"required"`
		LogFormat   string `json:"log_format"`
	}
	var req CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if errors := pkg.ValidateAndFormat(&req); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Validation failed", errors))
	}

	if err := h.projectService.CreateProject(req.Name, req.PathLog, req.ProjectType, req.LogFormat); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to create project", nil))
	}
	return c.Status(fiber.StatusCreated).JSON(pkg.GenerateResponse("Project created successfully", nil))
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	type UpdateProjectRequest struct {
		Name        string `json:"name" validate:"required,min=2"`
		PathLog     string `json:"path_log" validate:"required"`
		ProjectType string `json:"project_type" validate:"required"`
		LogFormat   string `json:"log_format"`
	}
	var req UpdateProjectRequest
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

	project := &model.Project{
		ID:          uint(id),
		Name:        req.Name,
		PathLog:     req.PathLog,
		ProjectType: req.ProjectType,
		LogFormat:   req.LogFormat,
	}
	if err := h.projectService.UpdateProject(project); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update project", nil))
	}
	return c.JSON(pkg.GenerateResponse("Project updated successfully", project))
}

func (h *ProjectHandler) GetProjectLogs(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid project ID", nil))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "100"))
	offset := (page - 1) * limit

	logs, total, err := h.logService.GetLogsByProjectID(uint(id), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve logs", nil))
	}

	return c.JSON(pkg.GenerateResponse("Logs retrieved successfully", fiber.Map{
		"data":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	}))
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid param ID", nil))
	}
	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to delete project", nil))
	}
	return c.JSON(pkg.GenerateResponse("Project deleted successfully", nil))
}

func (h *ProjectHandler) SetupProjectRoutes(router fiber.Router) {
	projectGroup := router.Group("/projects")
	projectGroup.Post("/", h.CreateProject)
	projectGroup.Get("/", h.GetAllProjects)
	projectGroup.Get("/:id", h.GetProjectByID)
	projectGroup.Put("/:id", h.UpdateProject)
	projectGroup.Delete("/:id", h.DeleteProject)
	projectGroup.Get("/:id/logs", h.GetProjectLogs)
}
