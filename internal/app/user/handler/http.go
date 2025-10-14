package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
)

type UserHandler struct {
	service domain.UserService
	log     logger.Logger
}

func NewUserHandler(srv *gin.RouterGroup, service domain.UserService, log logger.Logger) *UserHandler {
	handlers := &UserHandler{
		log:     log,
		service: service,
	}

	srv.GET("/test", handlers.test)

	return handlers
}

func (h *UserHandler) test(c *gin.Context) {
	userID, ok := c.Get(domain.CtxUserIDKey)
	if !ok {
		c.JSON(400, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrInternalServerError.Error(),
			},
		})
		return
	}
	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(400, domain.ErrorResponse{Errors: map[string]any{"error": err.Error()}})
		return
	}
	c.JSON(200, domain.SuccessResponse{
		Message: "success",
		Data:    uid,
	})
}
