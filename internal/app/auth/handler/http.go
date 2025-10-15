package handler

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/internal/server/middleware"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/ilhamgepe/gepay/pkg/utils"
)

type AuthHandler struct {
	authService domain.AuthService
	mw          *middleware.Middleware
	log         logger.Logger
}

func NewAuthHandler(srv *gin.RouterGroup, authService domain.AuthService, mw *middleware.Middleware, log logger.Logger) *AuthHandler {
	handlers := &AuthHandler{
		authService: authService,
		mw:          mw,
		log:         log,
	}

	srv.POST("/sign-in", handlers.signIn)
	srv.POST("/sign-up", handlers.signUp)
	srv.GET("/me", mw.WithAuth, handlers.me)
	srv.GET("/refresh", handlers.refresh)
	return handlers
}

func (h *AuthHandler) signUp(c *gin.Context) {
	var req domain.SignUpReq
	if err := c.ShouldBind(&req); err != nil {
		result := utils.GenerateMessage(err, reflect.TypeOf(req))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "invalid request",
			Errors:  result,
		})
		return
	}

	req.Role = "merchant"
	err := h.authService.SignupLocal(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
	})
}

func (h *AuthHandler) signIn(c *gin.Context) {
	var req domain.SignInReq
	if err := c.ShouldBind(&req); err != nil {
		result := utils.GenerateMessage(err, reflect.TypeOf(req))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "invalid request",
			Errors:  result,
		})
		return
	}

	res, err := h.authService.Signin(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
		Data:    res,
	})
}

func (h *AuthHandler) me(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	res, err := h.authService.Me(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
		Data:    res,
	})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if bearerToken == nil || len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrUnauthorized.Error(),
			},
		})
		return
	}

	// validate token
	res, err := h.authService.Refresh(c.Request.Context(), bearerToken[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
		Data:    res,
	})
}
