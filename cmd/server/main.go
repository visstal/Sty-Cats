package main

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "spy-cat-agency/docs"
	"spy-cat-agency/internal/api/http/handlers"
	custommw "spy-cat-agency/internal/api/http/middleware"
	"spy-cat-agency/internal/api/http/routes"
	"spy-cat-agency/internal/application/services"
	"spy-cat-agency/internal/infrastructure/database"
	"spy-cat-agency/internal/infrastructure/mock_data"
	"spy-cat-agency/internal/infrastructure/repositories"
	"spy-cat-agency/pkg/validator"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title Spy Cat Agency API
// @version 1.0
// @description A simple API for managing spy cats, missions, and targets
// @host localhost:3001
// @BasePath /
func main() {
	_ = godotenv.Load()

	dbConfig := database.Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvInt("DB_PORT", 5432),
		Database:        getEnv("DB_NAME", "spy_cats"),
		Username:        getEnv("DB_USER", "spy_user"),
		Password:        getEnv("DB_PASSWORD", "spy_password"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 5)) * time.Minute,
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run auto-migration: %v", err)
	}

	mockService := mock_data.NewMockDataService(db)
	if err := mockService.WipeAndSeedData(); err != nil {
		log.Fatalf("Failed to initialize mock data: %v", err)
	}
	log.Println("Mock data initialized successfully")

	catRepo := repositories.NewCatRepository(db)
	missionRepo := repositories.NewMissionRepository(db.DB).(*repositories.MissionRepository)
	targetRepo := repositories.NewTargetRepository(db.DB).(*repositories.TargetRepository)

	missionService := services.NewMissionService(db, missionRepo, targetRepo, catRepo.(*repositories.CatRepository))

	catHandler := handlers.NewCatHandler(catRepo)
	missionHandler := handlers.NewMissionHandler(missionService)

	e := echo.New()

	e.Validator = validator.NewValidator()

	e.Use(custommw.LoggingMiddleware())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	routes.SetupRoutes(e, catHandler, missionHandler)

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
