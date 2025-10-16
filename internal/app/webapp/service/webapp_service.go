package service

import (
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/internal/server/security"
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type webappService struct {
	userService     domain.UserService
	merchantService domain.MerchantService
	log             logger.Logger
	security        *security.Security
	db              *sqlx.DB
	config          config.App
}

func NewWebappService(userService domain.UserService, merchantService domain.MerchantService, log logger.Logger, security *security.Security, db *sqlx.DB, config config.App) domain.WebappService {
	return &webappService{
		userService:     userService,
		merchantService: merchantService,
		log:             log,
		security:        security,
		db:              db,
		config:          config,
	}
}
