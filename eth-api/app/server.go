package app

import (
	"eth-api/app/handlers"
	"eth-api/app/middleware"
	"eth-api/app/repositories"
	"eth-api/app/routes"
	"eth-api/app/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewApp() *fiber.App {
	db := InitDatabase()
	handler := InitEthBlockHandler(db)

	app := fiber.New()

	// Apply logging middleware
	app.Use(middleware.LoggingMiddleware)

	// Set up routes
	routes.SetupRoutes(app, handler)

	return app
}

func InitDatabase() *gorm.DB {

	_ = godotenv.Load(".env")

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")
	connectTimeout := os.Getenv("POSTGRES_CONNECT_TIMEOUT")

	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s connect_timeout=%s",
		host, port, user, password, dbName, sslMode, timezone, connectTimeout)

	db, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	//// Enable logging mode
	////db.LogMode(true)
	//
	//// Migrate the database to create the EthBlock table
	//err = db.AutoMigrate(&models.EthBlock{})
	//if err != nil {
	//	log.Fatalf("failed to migrate database: %s", err)
	//}

	return db
}

func InitEthBlockHandler(db *gorm.DB) *handlers.EthBlockHandler {
	ethBlockRepo := repositories.NewEthBlockRepository(db)
	ethBlockService := services.NewEthBlockService(ethBlockRepo)
	ethBlockHandler := handlers.NewEthBlockHandler(ethBlockService)
	return ethBlockHandler
}
