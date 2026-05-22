package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/mohammed-ayoub-js/endline-backend/internal/config"
	handlers "github.com/mohammed-ayoub-js/endline-backend/internal/handler"
	"github.com/mohammed-ayoub-js/endline-backend/internal/middleware"
	"github.com/mohammed-ayoub-js/endline-backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")

	err = db.AutoMigrate(&models.User{}, &models.RefreshToken{}, &models.Subject{}, &models.Video{}, &models.UserVideoProgress{})
	if err != nil {
		log.Fatal("AutoMigration failed:", err)
	}
	log.Println("Database migrated successfully")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://yourfrontend.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Use(middleware.SecurityHeaders())
	app.Use(middleware.RateLimiter())

	api := app.Group("/api")
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register(db, cfg))
	auth.Post("/login", handlers.Login(db, cfg))
	auth.Post("/refresh", handlers.Refresh(db, cfg))
	auth.Post("/logout", handlers.Logout(db))

	protected := api.Group("/protected", middleware.JWTAuth(cfg))
	protected.Get("/profile", handlers.Profile)
    protected.Get("/videos/search", handlers.SearchVideos)
	protected.Post("/videos", handlers.SaveVideoToDB(db))
	protected.Get("/library", handlers.GetLibraryVideos(db))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}