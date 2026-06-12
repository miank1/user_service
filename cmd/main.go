package main

import (
	"log"
	"os"
	"user-service/internal/handler"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/miank1/ecommerce_backend/pkg/db"
	logger "github.com/miank1/ecommerce_backend/pkg/logger"
	"github.com/miank1/ecommerce_backend/pkg/middleware"
)

func main() {

	logger.Init()
	defer logger.Sync()

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ No .env file found")
	}
	log.Println("Loaded DSN:", os.Getenv("DATABASE_DSN"))

	dsn := os.Getenv("DATABASE_DSN")

	gormDB, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// Auto migrate User model
	if err := gormDB.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("❌ Auto migrate failed: %v", err)
	} else {
		log.Println("✅ User table migration successful!")
	}

	// Initialize repository, service, handler
	repo := repository.NewUserRepository(gormDB)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	// Initialize Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "userservice is up"})
	})

	// Database connectivity check
	r.GET("/health/db", func(c *gin.Context) {
		sqlDB, err := gormDB.DB()
		if err != nil {
			c.JSON(500, gin.H{"db": "error", "details": err.Error()})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"db": "not reachable", "details": err.Error()})
			return
		}
		c.JSON(200, gin.H{"db": "connected ✅"})
	})

	api := r.Group("/users")

	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", h.Me)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("✅ USER SERVICE running on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

}
