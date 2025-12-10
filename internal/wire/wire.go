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

// InitializeApp initializes the application with all dependencies
func InitializeApp(cfg *config.Config) (*handler.UserHandler, error) {
	wire.Build(
		// Database
		provideDatabase,
		// Redis
		provideRedis,
		// Repository
		repository.NewUserRepository,
		// Service
		service.NewUserService,
		// Handler
		handler.NewUserHandler,
	)
	return nil, nil
}

func provideDatabase(cfg *config.Config) (*gorm.DB, error) {
	return database.NewMySQL(&cfg.Database)
}

func provideRedis(cfg *config.Config) (*redis.Client, error) {
	return pkgredis.NewRedis(&cfg.Redis)
}
