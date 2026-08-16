package http_handler

import (
	"errors"
	"strconv"

	service "github.com/gabutlabs/godevopin/internal/services"
	"github.com/gabutlabs/godevopin/pkg"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MySQLActivityHandler struct{ service service.MySQLActivityService }

func NewMySQLActivityHandler(svc service.MySQLActivityService) *MySQLActivityHandler {
	return &MySQLActivityHandler{service: svc}
}

func (h *MySQLActivityHandler) ListTargets(c *fiber.Ctx) error {
	targets, err := h.service.ListMySQLTargets()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve MySQL targets", nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL targets retrieved successfully", targets))
}

func (h *MySQLActivityHandler) CreateTarget(c *fiber.Ctx) error {
	var request service.MySQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	target, err := h.service.CreateMySQLTarget(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.Status(fiber.StatusCreated).JSON(pkg.GenerateResponse("MySQL target created successfully", target))
}

func (h *MySQLActivityHandler) UpdateTarget(c *fiber.Ctx) error {
	id, err := parseMySQLTargetID(c)
	if err != nil {
		return err
	}
	var request service.MySQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	target, err := h.service.UpdateMySQLTarget(id, request)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("MySQL target not found", nil))
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL target updated successfully", target))
}

func (h *MySQLActivityHandler) DeleteTarget(c *fiber.Ctx) error {
	id, err := parseMySQLTargetID(c)
	if err != nil {
		return err
	}
	if err := h.service.DeleteMySQLTarget(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to delete MySQL target", nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL target deleted successfully", nil))
}

func (h *MySQLActivityHandler) TestTarget(c *fiber.Ctx) error {
	var request service.MySQLTargetRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}
	version, err := h.service.TestMySQLTarget(request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL connection successful", fiber.Map{"server_version": version}))
}

func (h *MySQLActivityHandler) GetLiveActivity(c *fiber.Ctx) error {
	targetID, err := parseMySQLOptionalUint(c.Query("target_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid target_id", nil))
	}
	limit, err := parseMySQLLimit(c.Query("limit"), 500)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid activity limit", nil))
	}
	result, err := h.service.GetLiveMySQLActivity(targetID, c.Query("search"), limit)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(pkg.GenerateErrorResponse("MySQL target not found", nil))
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve MySQL activity", nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL activity retrieved successfully", result))
}

func (h *MySQLActivityHandler) GetHistory(c *fiber.Ctx) error {
	targetID, err := parseMySQLOptionalUint(c.Query("target_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid target_id", nil))
	}
	minimumDuration, err := strconv.ParseInt(c.Query("min_duration_ms", "0"), 10, 64)
	if err != nil || minimumDuration < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid min_duration_ms", nil))
	}
	limit, err := parseMySQLLimit(c.Query("limit"), 500)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid activity history limit", nil))
	}
	history, err := h.service.GetMySQLActivityHistory(c.Query("activity_key"), targetID, c.Query("filter", "1h"), c.Query("search"), c.Query("state"), minimumDuration, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse(err.Error(), nil))
	}
	return c.JSON(pkg.GenerateResponse("MySQL activity history retrieved successfully", history))
}

func (h *MySQLActivityHandler) SetupRoutes(router fiber.Router) {
	group := router.Group("/mysql-activity")
	group.Get("/targets", h.ListTargets)
	group.Post("/targets", h.CreateTarget)
	group.Put("/targets/:id", h.UpdateTarget)
	group.Delete("/targets/:id", h.DeleteTarget)
	group.Post("/targets/test", h.TestTarget)
	group.Get("/live", h.GetLiveActivity)
	group.Get("/history", h.GetHistory)
}

func parseMySQLTargetID(c *fiber.Ctx) (uint, error) {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid MySQL target ID", nil))
	}
	return uint(id), nil
}

func parseMySQLOptionalUint(value string) (uint, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	return uint(parsed), err
}

func parseMySQLLimit(value string, fallback int) (int, error) {
	if value == "" {
		return fallback, nil
	}
	limit, err := strconv.Atoi(value)
	if err != nil || limit < 1 {
		return 0, errors.New("invalid limit")
	}
	return limit, nil
}
