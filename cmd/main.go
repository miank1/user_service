package main

import (
	"fmt"
	"log"
	"os"
	"user-service/internal/handler"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/db"
	"user-service/pkg/logger"
	"user-service/pkg/middleware"

	_ "github.com/miank1/user_service/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User Service API
// @version 1.0
// @description User Management Service
// @host localhost:8081
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

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

	// Public routes
	// Register godoc
	//
	//	@Summary		Register user
	//	@Description	Register a new user
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			request	body		dto.RegisterRequest	true	"Register Request"
	//	@Success		201		{object}	dto.UserResponse
	//	@Failure		400		{object}	map[string]string
	//	@Router			/users/register [post]
	api.POST("/register", h.Register)

	// Login godoc
	//
	//	@Summary		Login user
	//	@Description	Authenticate user and return JWT
	//	@Tags			Users
	//	@Accept			json
	//	@Produce		json
	//	@Param			request	body		dto.LoginRequest	true	"Login Request"
	//	@Success		200		{object}	dto.LoginResponse
	//	@Failure		401		{object}	map[string]string
	//	@Router			/users/login [post]
	api.POST("/login", h.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTAuth())

	protected.GET("/me", h.Me)

	// Me godoc
	//
	//	@Summary		Get current user
	//	@Description	Get logged in user profile
	//	@Tags			Users
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{object}	dto.UserResponse
	//	@Router			/users/me [get]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Start server
	fmt.Printf("✅ ***** USER SERVICE ***** running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
