package main

import (
	"fmt"
	"log"

	"github.com/IndigoCloud6/go-web-template/internal/config"
	"github.com/IndigoCloud6/go-web-template/internal/middleware"
	"github.com/IndigoCloud6/go-web-template/internal/model"
	"github.com/IndigoCloud6/go-web-template/internal/wire"
	"github.com/IndigoCloud6/go-web-template/pkg/database"
	"github.com/IndigoCloud6/go-web-template/pkg/logger"
	"github.com/IndigoCloud6/go-web-template/pkg/response"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/IndigoCloud6/go-web-template/docs"
)

// @title Go Web Template API
// @version 1.0
// @description A production-ready Go web service template
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(&cfg.Logger); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database and auto-migrate
	db, err := database.NewMySQL(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database")
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
		logger.Fatal("Failed to auto-migrate database")
	}

	// Initialize app with Wire
	handlers, err := wire.InitializeApp(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize app")
	}

	// Create Gin router
	r := gin.New()

	// Apply middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, map[string]interface{}{
			"status": "ok",
		})
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", handlers.UserHandler.CreateUser)
			users.GET("", handlers.UserHandler.ListUsers)
			users.GET("/:id", handlers.UserHandler.GetUser)
			users.PUT("/:id", handlers.UserHandler.UpdateUser)
			users.DELETE("/:id", handlers.UserHandler.DeleteUser)
		}

		products := v1.Group("/products")
		{
			products.POST("", handlers.ProductHandler.CreateProduct)
			products.GET("", handlers.ProductHandler.ListProducts)
			products.GET("/:id", handlers.ProductHandler.GetProduct)
			products.PUT("/:id", handlers.ProductHandler.UpdateProduct)
			products.DELETE("/:id", handlers.ProductHandler.DeleteProduct)
		}
	}

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info(fmt.Sprintf("Server starting on %s", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
