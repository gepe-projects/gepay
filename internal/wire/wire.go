package wire

import (
	authHandler "github.com/ilhamgepe/gepay/internal/app/auth/handler"
	authService "github.com/ilhamgepe/gepay/internal/app/auth/service"
	merchantHandler "github.com/ilhamgepe/gepay/internal/app/merchant/handler"
	merchantRepo "github.com/ilhamgepe/gepay/internal/app/merchant/repository"
	merchantService "github.com/ilhamgepe/gepay/internal/app/merchant/service"
	userHandler "github.com/ilhamgepe/gepay/internal/app/user/handler"
	userRepo "github.com/ilhamgepe/gepay/internal/app/user/repository"
	userService "github.com/ilhamgepe/gepay/internal/app/user/service"
	"github.com/ilhamgepe/gepay/internal/server"
	"github.com/ilhamgepe/gepay/internal/server/middleware"
	"github.com/ilhamgepe/gepay/internal/server/security"
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/database"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/ilhamgepe/gepay/pkg/redis"
)

type Apps struct {
	Server *server.Server
}

func InitializeApps() *Apps {
	config := provideConfig()
	log := provideLogger(config)

	// db
	db := database.InitDB(&config.Database, log)

	// redis
	rdb := redis.NewClient(config.Redis, log)

	// repo
	userRepo := userRepo.NewUserRepository(db, log)
	merchantRepo := merchantRepo.NewMerchantRepo(db, log)

	// security
	security := security.NewSecurity(config, rdb, log)

	// service
	merchantService := merchantService.NewMerchantService(merchantRepo, db, log)
	userService := userService.NewUserService(userRepo, log)
	authService := authService.NewAuthService(userService, merchantService, log, security, db, config)

	// server
	server := server.NewServer(&config.Server, log)
	v1 := server.App.Group("/api/v1")

	// middleware
	mw := middleware.NewMiddlewares(rdb, security, config, log)

	// handler
	userHandler.NewUserHandler(v1.Group("/users", mw.WithAuth), userService, log)
	authHandler.NewAuthHandler(v1.Group("/auth"), authService, mw, log)
	merchantHandler.NewMerchantHandler(v1.Group("merchants"), merchantService, mw, log)

	return &Apps{
		Server: server,
	}
}

func provideConfig() config.App {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	return cfg
}

func provideLogger(cfg config.App) logger.Logger {
	return logger.New("gepay", cfg)
}
