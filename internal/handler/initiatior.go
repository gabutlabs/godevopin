package handler

import (
	"github.com/gabutlabs/godevopin/internal/config"
	"github.com/gabutlabs/godevopin/internal/database"
	http_handler "github.com/gabutlabs/godevopin/internal/handler/http"
	"github.com/gabutlabs/godevopin/internal/handler/socket"
	"github.com/gabutlabs/godevopin/internal/repository"
	service "github.com/gabutlabs/godevopin/internal/services"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupHandlers(router fiber.Router, dbs *database.Connections, config *config.Config) {
	userRepo := repository.NewUserRepository(dbs.App)
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
	systemMetricRepo := repository.NewSystemMetricRepository(dbs.Metrics)
	systemMetricService := service.NewSystemMetricService(systemMetricRepo)
	systemMetricHandler := http_handler.NewSystemMetricHandler(systemMetricService)
	systemMetricHandler.SetupSystemMetricRoutes(router)

	// Setup process monitoring routes
	processMetricRepo := repository.NewProcessMetricRepository(dbs.Metrics)
	processMonitoringService := service.NewProcessMonitoringService(processMetricRepo)
	processHandler := http_handler.NewProcessMonitoringHandler(processMonitoringService)
	processHandler.SetupProcessMonitoringRoutes(router)

	// Setup PostgreSQL activity routes. Target credentials are stored encrypted
	// with the application secret and are never included in API responses.
	postgresTargetRepo := repository.NewPostgreSQLTargetRepository(dbs.App)
	postgresActivityRepo := repository.NewPostgreSQLActivityRepository(dbs.Metrics)
	postgresActivityService := service.NewPostgreSQLActivityService(postgresTargetRepo, postgresActivityRepo, config.AppSetting.JWTSecret)
	postgresActivityHandler := http_handler.NewPostgreSQLActivityHandler(postgresActivityService)
	postgresActivityHandler.SetupRoutes(router)

	// Setup MySQL activity routes.
	mysqlTargetRepo := repository.NewMySQLTargetRepository(dbs.App)
	mysqlActivityRepo := repository.NewMySQLActivityRepository(dbs.Metrics)
	mysqlActivityService := service.NewMySQLActivityService(mysqlTargetRepo, mysqlActivityRepo, config.AppSetting.JWTSecret)
	mysqlActivityHandler := http_handler.NewMySQLActivityHandler(mysqlActivityService)
	mysqlActivityHandler.SetupRoutes(router)

	// Setup WorkerService routes
	workerServiceRepo := repository.NewWorkerServiceRepository(dbs.App)
	workerServiceService := service.NewWorkerServiceService(workerServiceRepo)
	workerServiceHandler := http_handler.NewWorkerServiceHandler(workerServiceService)
	workerServiceHandler.SetupWorkerServiceRoutes(router)

	// Setup Alarm routes
	alarmRepo := repository.NewAlarmRepository(dbs.App)
	alarmService := service.NewAlarmService(alarmRepo)
	alarmHandler := http_handler.NewAlarmHandler(alarmService)
	alarmHandler.SetupAlarmRoutes(router)

	// Setup Docker routes
	dockerService := service.NewDockerService()
	dockerHandler := http_handler.NewDockerHandler(dockerService)
	dockerHandler.SetupDockerRoutes(router)

	// Setup Project routes
	projectRepo := repository.NewProjectRepository(dbs.App)
	projectService := service.NewProjectService(projectRepo)
	logHistoryRepo := repository.NewLogHistoryRepository(dbs.Logs)
	logHistoryService := service.NewLogHistoryService(logHistoryRepo)
	projectHandler := http_handler.NewProjectHandler(projectService, logHistoryService)
	projectHandler.SetupProjectRoutes(router)

	// Setup Setting routes
	settingRepo := repository.NewSettingRepository(dbs.App)
	settingService := service.NewSettingService(settingRepo)
	settingHandler := http_handler.NewSettingHandler(settingService)
	settingHandler.SetupSettingRoutes(router)
}

func SetupSocketHandlers(wsGroup fiber.Router, dbs *database.Connections, config *config.Config) {
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
