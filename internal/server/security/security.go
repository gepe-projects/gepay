package security

import (
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Security struct {
	config config.App
	rdb    *redis.Client
	log    logger.Logger
}

func NewSecurity(config config.App, rdb *redis.Client, log logger.Logger) *Security {
	return &Security{config: config, rdb: rdb, log: log}
}
