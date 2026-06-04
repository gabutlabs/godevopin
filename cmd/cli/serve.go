// cmd/nama-aplikasi/cmd/serve.go
package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/gabutlabs/godevopin/internal/config"
	"github.com/gabutlabs/godevopin/internal/database"
	"github.com/gabutlabs/godevopin/internal/handler"
	"github.com/gabutlabs/godevopin/internal/model"
	"github.com/gabutlabs/godevopin/internal/web"
	"github.com/gabutlabs/godevopin/pkg"
	"github.com/gofiber/fiber/v2" // Impor Fiber
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/cobra"
)

// Variabel untuk menyimpan nilai dari flag --port
var port int

// serveCmd merepresentasikan command 'serve'
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Menjalankan server web Fiber",
	Long:  `Menjalankan server web HTTP yang dibangun menggunakan framework Fiber.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Inisialisasi aplikasi Fiber
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Fatalf("could not load config: %v", err)
		}

		// 2. Hubungkan ke database menggunakan GORM
		db, err := database.ConnectDB(cfg.Database) // Sekarang db bertipe *gorm.DB
		if err != nil {
			log.Fatalf("could not connect to database: %v", err)
		}

		err = db.AutoMigrate(
			&model.User{},
			&model.SystemMetric{},
			&model.WorkerService{},
			&model.ActiveAlarm{},
			&model.AlarmHistory{},
			&model.Project{},
			&model.AppSetting{},
			&model.LogHistory{},
		)

		if err != nil {
			log.Fatalf("could not migrate database: %v", err)
		}

		// Seeder untuk AppSetting
		var count int64
		db.Model(&model.AppSetting{}).Count(&count)
		if count == 0 {
			defaultSetting := model.AppSetting{
				ID:                            1,
				MonitoringIntervalSeconds:     10,
				AlarmCheckIntervalSeconds:     60,
				AlarmRepeatIntervalMinutes:    15,
				SystemCPUCriticalPercent:      90.0,
				SystemDiskCriticalPercent:     85.0,
				SystemMemCriticalPercent:      85.0,
				WorkerHeartbeatTimeoutSeconds: 300,
			}
			if err := db.Create(&defaultSetting).Error; err != nil {
				log.Printf("could not seed default settings: %v", err)
			} else {
				fmt.Println("Default settings seeded successfully.")
			}
		}

		// Seeder untuk User Admin
		var userCount int64
		db.Model(&model.User{}).Count(&userCount)
		if userCount == 0 {
			hashedPassword, err := pkg.HashPassword("password1!")
			if err != nil {
				log.Printf("could not hash default admin password: %v", err)
			} else {
				defaultAdmin := model.User{
					Name:     "Admin",
					Email:    "admin@gabutngoding.com",
					Password: hashedPassword,
				}
				if err := db.Create(&defaultAdmin).Error; err != nil {
					log.Printf("could not seed default admin user: %v", err)
				} else {
					fmt.Println("Default admin user seeded successfully.")
				}
			}
		}
		app := fiber.New()
		app.Use(cors.New())
		// -- HTTP Routes --
		// Serve embedded frontend files
		app.Get("/", func(c *fiber.Ctx) error {
			indexHtml, err := web.IndexHtml.ReadFile("dist/index.html")
			if err != nil {
				return c.Status(404).SendString("index.html not found")
			}
			c.Set("Content-Type", "text/html")
			return c.Send(indexHtml)
		})

		// Serve embedded assets (CSS, JS, fonts, etc.) from the embedded filesystem
		app.Get("/assets/*", func(c *fiber.Ctx) error {
			assetPath := "dist/assets/" + c.Params("*")
			assetData, err := web.Assets.ReadFile(assetPath)
			if err != nil {
				return c.Status(404).SendString("Asset not found")
			}

			// Set appropriate content type based on file extension
			filePath := c.Params("*")
			switch {
			case strings.HasSuffix(filePath, ".css"):
				c.Set("Content-Type", "text/css")
			case strings.HasSuffix(filePath, ".js"):
				c.Set("Content-Type", "application/javascript")
			case strings.HasSuffix(filePath, ".json"):
				c.Set("Content-Type", "application/json")
			case strings.HasSuffix(filePath, ".png"):
				c.Set("Content-Type", "image/png")
			case strings.HasSuffix(filePath, ".jpg") || strings.HasSuffix(filePath, ".jpeg"):
				c.Set("Content-Type", "image/jpeg")
			case strings.HasSuffix(filePath, ".gif"):
				c.Set("Content-Type", "image/gif")
			case strings.HasSuffix(filePath, ".ico"):
				c.Set("Content-Type", "image/x-icon")
			case strings.HasSuffix(filePath, ".woff"):
				c.Set("Content-Type", "font/woff")
			case strings.HasSuffix(filePath, ".woff2"):
				c.Set("Content-Type", "font/woff2")
			case strings.HasSuffix(filePath, ".ttf"):
				c.Set("Content-Type", "font/ttf")
			case strings.HasSuffix(filePath, ".eot"):
				c.Set("Content-Type", "application/vnd.ms-fontobject")
			default:
				c.Set("Content-Type", "application/octet-stream")
			}

			return c.Send(assetData)
		})

		// Serve favicon
		app.Get("/favicon.ico", func(c *fiber.Ctx) error {
			favicon, err := web.Favicon.ReadFile("dist/favicon.ico")
			if err != nil {
				return c.Status(404).SendString("favicon.ico not found")
			}
			c.Set("Content-Type", "image/x-icon")
			return c.Send(favicon)
		})
		apiGroup := app.Group("/api")
		handler.SetupHandlers(apiGroup, db, &cfg)
		// -- END HTTP Routes --

		// -- Socket Routes --
		wsGroup := app.Group("/ws")
		handler.SetupSocketHandlers(wsGroup, db, &cfg)
		// -- END Socket Routes --

		// Mulai server pada port yang ditentukan
		fmt.Printf("Server Fiber berjalan di http://localhost:%d\n", port)
		log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
	},
}

func init() {
	// Menambahkan flag --port atau -p ke subcommand 'serve'
	serveCmd.Flags().IntVarP(&port, "port", "p", 8000, "Port untuk menjalankan server web")
}
