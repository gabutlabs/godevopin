package http_handler

import (
	"strconv"

	"github.com/gabutlabs/devopin/internal/model"
	service "github.com/gabutlabs/devopin/internal/services"
	"github.com/gabutlabs/devopin/pkg"
	"github.com/gofiber/fiber/v2"
)

// AlarmHandler follows the same pattern as other handlers in the project
type AlarmHandler struct {
	alarmService service.AlarmService
}

// Request structs for the handlers
type UpdateActiveAlarmRequest struct {
	Status         string `json:"status" validate:"required,oneof=FIRING ACKNOWLEDGED"`
	AcknowledgedBy string `json:"acknowledged_by"`
	Message        string `json:"message"`
}

type AcknowledgeAlarmRequest struct {
	AcknowledgedBy string `json:"acknowledged_by" validate:"required"`
}

func NewAlarmHandler(alarmService service.AlarmService) *AlarmHandler {
	return &AlarmHandler{
		alarmService: alarmService,
	}
}

// GetAllActiveAlarms godoc
// @Summary Get all active alarms
// @Description Get list of all active alarms
// @Tags alarms
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /alarms/active [get]
func (h *AlarmHandler) GetAllActiveAlarms(c *fiber.Ctx) error {
	alarms, err := h.alarmService.GetAllActiveAlarms()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve active alarms", nil))
	}

	return c.JSON(pkg.GenerateResponse("Active alarms retrieved successfully", alarms))
}

// GetAllAlarmHistoryPaginated godoc
// @Summary Get alarm history with pagination
// @Description Get paginated list of alarm history
// @Tags alarms
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param sort query string false "Sort order" default(created_at desc)
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /alarms/history [get]
func (h *AlarmHandler) GetAllAlarmHistoryPaginated(c *fiber.Ctx) error {

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	sort := c.Query("sort", "created_at DESC")
	filterStatus := c.Query("status", "")
	result, err := h.alarmService.GetAllAlarmHistoryPaginated(page, limit, sort, filterStatus)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve alarm history", nil))
	}

	return c.JSON(pkg.GenerateResponse("Alarm history retrieved successfully", result))
}

// GetActiveAlarmByStatus godoc
// @Summary Get active alarms by status
// @Description Get list of active alarms filtered by status
// @Tags alarms
// @Accept json
// @Produce json
// @Param status path string true "Alarm status (FIRING or ACKNOWLEDGED)"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /alarms/active/status/{status} [get]
func (h *AlarmHandler) GetActiveAlarmByStatus(c *fiber.Ctx) error {
	statusStr := c.Params("status")
	if statusStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Status parameter is required", nil))
	}

	var status model.ActiveAlarmStatus
	switch statusStr {
	case "FIRING":
		status = model.StatusFiring
	case "ACKNOWLEDGED":
		status = model.StatusAcknowledged
	default:
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid status, must be FIRING or ACKNOWLEDGED", nil))
	}

	alarms, err := h.alarmService.GetActiveAlarmsByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to retrieve active alarms by status", nil))
	}

	return c.JSON(pkg.GenerateResponse("Active alarms retrieved successfully", alarms))
}

// UpdateActiveAlarm godoc
// @Summary Update an active alarm
// @Description Update status and other details of an active alarm
// @Tags alarms
// @Accept json
// @Produce json
// @Param alarmName path string true "Alarm name"
// @Param target path string true "Alarm target"
// @Param request body UpdateActiveAlarmRequest true "Update alarm request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /alarms/active/{alarmName}/{target} [put]
func (h *AlarmHandler) UpdateActiveAlarm(c *fiber.Ctx) error {
	alarmName := c.Params("alarmName")
	target := c.Params("target")

	if alarmName == "" || target == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Alarm name and target are required", nil))
	}

	var req UpdateActiveAlarmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}

	var status model.ActiveAlarmStatus
	switch req.Status {
	case "FIRING":
		status = model.StatusFiring
	case "ACKNOWLEDGED":
		status = model.StatusAcknowledged
	default:
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Invalid status, must be FIRING or ACKNOWLEDGED", nil))
	}

	var acknowledgedBy *string
	if req.AcknowledgedBy != "" {
		acknowledgedBy = &req.AcknowledgedBy
	}

	err := h.alarmService.UpdateActiveAlarm(alarmName, target, status, acknowledgedBy, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to update active alarm", nil))
	}

	return c.JSON(pkg.GenerateResponse("Active alarm updated successfully", nil))
}

// AcknowledgeActiveAlarm godoc
// @Summary Acknowledge an active alarm
// @Description Mark an active alarm as acknowledged
// @Tags alarms
// @Accept json
// @Produce json
// @Param alarmName path string true "Alarm name"
// @Param target path string true "Alarm target"
// @Param request body AcknowledgeAlarmRequest true "Acknowledge alarm request"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /alarms/active/acknowledge/{alarmName}/{target} [post]
func (h *AlarmHandler) AcknowledgeActiveAlarm(c *fiber.Ctx) error {
	alarmName := c.Params("alarmName")
	target := c.Params("target")

	if alarmName == "" || target == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Alarm name and target are required", nil))
	}

	var req AcknowledgeAlarmRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("Cannot parse JSON", nil))
	}

	if req.AcknowledgedBy == "" {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.GenerateErrorResponse("AcknowledgedBy is required", nil))
	}

	err := h.alarmService.AcknowledgeActiveAlarm(alarmName, target, req.AcknowledgedBy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.GenerateErrorResponse("Failed to acknowledge active alarm", nil))
	}

	return c.JSON(pkg.GenerateResponse("Active alarm acknowledged successfully", nil))
}

// SetupAlarmRoutes registers the alarm routes to the provided router
func (h *AlarmHandler) SetupAlarmRoutes(router fiber.Router) {
	alarmGroup := router.Group("/alarms")
	alarmGroup.Get("/active", h.GetAllActiveAlarms)
	alarmGroup.Get("/history", h.GetAllAlarmHistoryPaginated)
	alarmGroup.Get("/active/status/:status", h.GetActiveAlarmByStatus)
	alarmGroup.Put("/active/:alarmName/:target", h.UpdateActiveAlarm)
	alarmGroup.Post("/active/acknowledge/:alarmName/:target", h.AcknowledgeActiveAlarm)
}
