package handler

import (
	"github.com/gabutlabs/devopin/internal/config"
	http_handler "github.com/gabutlabs/devopin/internal/handler/http"
	"github.com/gabutlabs/devopin/internal/handler/socket"
	"github.com/gabutlabs/devopin/internal/repository"
	service "github.com/gabutlabs/devopin/internal/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupHandlers(router fiber.Router, db *gorm.DB, config *config.Config) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := http_handler.NewUserHandler(userService)
	authHandler := http_handler.NewAuthHandler(userService)
	authHandler.SetupAuthRoutes(router)
	router.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(config.AppSetting.JWTSecret)},
	}))
	router.Get("/auth/me", authHandler.Me)
	userHandler.SetupUserRoutes(router)

	// Setup SystemMetric routes
	systemMetricRepo := repository.NewSystemMetricRepository(db)
	systemMetricService := service.NewSystemMetricService(systemMetricRepo)
	systemMetricHandler := http_handler.NewSystemMetricHandler(systemMetricService)
	systemMetricHandler.SetupSystemMetricRoutes(router)

	// Setup WorkerService routes
	workerServiceRepo := repository.NewWorkerServiceRepository(db)
	workerServiceService := service.NewWorkerServiceService(workerServiceRepo)
	workerServiceHandler := http_handler.NewWorkerServiceHandler(workerServiceService)
	workerServiceHandler.SetupWorkerServiceRoutes(router)

	// Setup Alarm routes
	alarmRepo := repository.NewAlarmRepository(db)
	alarmService := service.NewAlarmService(alarmRepo)
	alarmHandler := http_handler.NewAlarmHandler(alarmService)
	alarmHandler.SetupAlarmRoutes(router)

	// Setup Docker routes
	dockerService := service.NewDockerService()
	dockerHandler := http_handler.NewDockerHandler(dockerService)
	dockerHandler.SetupDockerRoutes(router)

	// Setup Project routes
	projectRepo := repository.NewProjectRepository(db)
	projectService := service.NewProjectService(projectRepo)
	logHistoryRepo := repository.NewLogHistoryRepository(db)
	logHistoryService := service.NewLogHistoryService(logHistoryRepo)
	projectHandler := http_handler.NewProjectHandler(projectService, logHistoryService)
	projectHandler.SetupProjectRoutes(router)

	// Setup Setting routes
	settingRepo := repository.NewSettingRepository(db)
	settingService := service.NewSettingService(settingRepo)
	settingHandler := http_handler.NewSettingHandler(settingService)
	settingHandler.SetupSettingRoutes(router)
}

func SetupSocketHandlers(wsGroup fiber.Router, db *gorm.DB, config *config.Config) {
	// Setup Docker Socket routes
	wsGroup.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	dockerService := service.NewDockerService()
	dockerSocketHandler := socket.NewDockerSocketHandler(dockerService)
	dockerSocketHandler.SetupDockerSocketRoutes(wsGroup)
}
