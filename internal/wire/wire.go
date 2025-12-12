//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/IndigoCloud6/go-web-template/internal/config"
	"github.com/IndigoCloud6/go-web-template/internal/handler"
	"github.com/IndigoCloud6/go-web-template/internal/repository"
	"github.com/IndigoCloud6/go-web-template/internal/service"
	"github.com/IndigoCloud6/go-web-template/pkg/database"
	pkgredis "github.com/IndigoCloud6/go-web-template/pkg/redis"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Handlers holds all the application handlers
type Handlers struct {
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
	AuthHandler    *handler.AuthHandler
}

// InitializeApp initializes the application with all dependencies
func InitializeApp(cfg *config.Config) (*Handlers, error) {
	wire.Build(
		// Database
		provideDatabase,
		// Redis
		provideRedis,
		// JWT Config
		provideJWTConfig,
		// Repository
		repository.NewUserRepository,
		repository.NewProductRepository,
		// Service
		service.NewUserService,
		service.NewProductService,
		service.NewAuthService,
		// Handler
		handler.NewUserHandler,
		handler.NewProductHandler,
		handler.NewAuthHandler,
		// Handlers struct
		wire.Struct(new(Handlers), "*"),
	)
	return nil, nil
}

func provideDatabase(cfg *config.Config) (*gorm.DB, error) {
	return database.NewMySQL(&cfg.Database)
}

func provideRedis(cfg *config.Config) (*redis.Client, error) {
	return pkgredis.NewRedis(&cfg.Redis)
}

func provideJWTConfig(cfg *config.Config) *config.JWTConfig {
	return &cfg.JWT
}
