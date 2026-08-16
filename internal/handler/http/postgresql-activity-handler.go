package http_handler

import (
	"errors"
	"strconv"

	service "github.com/gabutlabs/godevopin/internal/services"
	"github.com/gabutlabs/godevopin/pkg"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PostgreSQLActivityHandler struct {
	service service.PostgreSQLActivityService
}

func NewPostgreSQLActivityHandler(svc service.PostgreSQLActivityService) *PostgreSQLActivityHandler {
	return &PostgreSQLActivityHandler{service: svc}
}

func (h *PostgreSQLActivityHandler) ListTargets(c *fiber.Ctx) error {
	targets, err := h.service.ListPostgreSQLTargets()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve PostgreSQL targets", nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL targets retrieved successfully", targets))
}

func (h *PostgreSQLActivityHandler) CreateTarget(c *fiber.Ctx) error {
	var request service.PostgreSQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	target, err := h.service.CreatePostgreSQLTarget(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.Status(fiber.StatusCreated).JSON(pkg.GenerateResponse("PostgreSQL target created successfully", target))
}

func (h *PostgreSQLActivityHandler) UpdateTarget(c *fiber.Ctx) error {
	id, err := parseTargetID(c)
	if err != nil {
		return err
	}
	var request service.PostgreSQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	target, err := h.service.UpdatePostgreSQLTarget(id, request)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("PostgreSQL target not found", nil))
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL target updated successfully", target))
}

func (h *PostgreSQLActivityHandler) DeleteTarget(c *fiber.Ctx) error {
	id, err := parseTargetID(c)
	if err != nil {
		return err
	}
	if err := h.service.DeletePostgreSQLTarget(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to delete PostgreSQL target", nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL target deleted successfully", nil))
}

func (h *PostgreSQLActivityHandler) TestTarget(c *fiber.Ctx) error {
	var request service.PostgreSQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	version, err := h.service.TestPostgreSQLTarget(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL connection successful", fiber.Map{"server_version": version}))
}

func (h *PostgreSQLActivityHandler) GetLiveActivity(c *fiber.Ctx) error {
	targetID, err := parseOptionalUint(c.Query("target_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid target_id", nil))
	}
	limit, err := parseLimit(c.Query("limit"), 500)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid activity limit", nil))
	}
	result, err := h.service.GetLivePostgreSQLActivity(targetID, c.Query("search"), limit)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("PostgreSQL target not found", nil))
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve PostgreSQL activity", nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL activity retrieved successfully", result))
}

func (h *PostgreSQLActivityHandler) GetHistory(c *fiber.Ctx) error {
	targetID, err := parseOptionalUint(c.Query("target_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid target_id", nil))
	}
	minimumDuration, err := strconv.ParseInt(c.Query("min_duration_ms", "0"), 10, 64)
	if err != nil || minimumDuration < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid min_duration_ms", nil))
	}
	blockedOnly, err := strconv.ParseBool(c.Query("blocked_only", "false"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid blocked_only", nil))
	}
	limit, err := parseLimit(c.Query("limit"), 500)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid activity history limit", nil))
	}
	history, err := h.service.GetPostgreSQLActivityHistory(
		c.Query("activity_key"),
		targetID,
		c.Query("filter", "1h"),
		c.Query("search"),
		c.Query("state"),
		blockedOnly,
		minimumDuration,
		limit,
	)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("PostgreSQL activity history retrieved successfully", history))
}

func (h *PostgreSQLActivityHandler) SetupRoutes(router fiber.Router) {
	group := router.Group("/postgresql-activity")
	group.Get("/targets", h.ListTargets)
	group.Post("/targets", h.CreateTarget)
	group.Put("/targets/:id", h.UpdateTarget)
	group.Delete("/targets/:id", h.DeleteTarget)
	group.Post("/targets/test", h.TestTarget)
	group.Get("/live", h.GetLiveActivity)
	group.Get("/history", h.GetHistory)
}

func parseTargetID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid PostgreSQL target ID", nil))
	}
	return uint(id), nil
}

func parseOptionalUint(value string) (uint, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(parsed), nil
}

func parseLimit(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 {
		return 0, errors.New("invalid limit")
	}
	return limit, nil
}
