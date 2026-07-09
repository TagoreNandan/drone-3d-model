package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"

	_ "gcs/docs"
	"gcs/internal/database"
	"gcs/internal/db"
	"gcs/internal/mockdata"
	"gcs/internal/modules/alerts"
	"gcs/internal/modules/missions"
	"gcs/internal/modules/telemetry"
	"gcs/internal/modules/vehicles"
)

var logger *zap.Logger

func initLogger() {
	var err error
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "production" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
}

// @title GCS API
// @version 0.1
// @description Ground Control Station backend API.
// @host localhost:8080
// @BasePath /
func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize logger
	initLogger()
	defer logger.Sync()

	logger.Info("Starting gcs server")

	// Initialize database pool
	pool, err := database.InitDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()
	logger.Info("Database pool connected successfully")

	queries := db.New(pool)

	// Get host from environment
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	// Get HTTP port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := host + ":" + port

	logger.Info("Starting HTTP server", zap.String("address", addr))

	// Create Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to gcs!",
		})
	})

	// Swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Register module routes under /api
	api := e.Group("/api")
	vehicles.NewController(vehicles.NewService(queries)).RegisterRoutes(api)
	missions.NewController(missions.NewService(queries)).RegisterRoutes(api)
	telemetry.NewController(telemetry.NewService(queries)).RegisterRoutes(api)
	alerts.NewController(alerts.NewService(queries)).RegisterRoutes(api)

	// Mock WebSocket route for Phase 2
	e.GET("/ws/mock/orientation", mockdata.HandleWS)

	// Start server
	if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
