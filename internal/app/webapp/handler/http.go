package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/internal/server/middleware"
	"github.com/ilhamgepe/gepay/pkg/logger"
)

type WebappHandler struct {
	services domain.WebappService
	mw       *middleware.Middleware
	log      logger.Logger
}

func NewWebappHandler(srv *gin.Engine, services domain.WebappService, mw *middleware.Middleware, log logger.Logger) {
	handlers := &WebappHandler{
		services: services,
		mw:       mw,
		log:      log,
	}
	authRoute := srv.Group("/auth")
	authRoute.POST("/sign-in", handlers.signIn)
	authRoute.POST("/sign-up", handlers.signUp)
	authRoute.GET("/me", mw.WithSessionAuth, handlers.me)
}
