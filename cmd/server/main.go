package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.AuthHandler.Login)
		}

		// Protected auth routes
		authProtected := v1.Group("/auth")
		authProtected.Use(middleware.JWTAuth(&cfg.JWT))
		{
			authProtected.POST("/refresh", handlers.AuthHandler.RefreshToken)
			authProtected.GET("/me", handlers.AuthHandler.GetCurrentUser)
		}

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

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server starting on %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Accept SIGINT (Ctrl+C) and SIGTERM (docker stop, k8s termination)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Set shutdown timeout from config, default to 30 seconds
	shutdownTimeout := cfg.Server.ShutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 30
	}

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimeout)*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server exited gracefully")
}
