package middleware

import (
	"github.com/ilhamgepe/gepay/internal/server/security"
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Middleware struct {
	rdb      *redis.Client
	security *security.Security
	config   config.App
	log      logger.Logger
}

func NewMiddlewares(rdb *redis.Client, security *security.Security, config config.App, log logger.Logger) *Middleware {
	return &Middleware{rdb: rdb, security: security, config: config, log: log}
}
