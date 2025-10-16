package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/utils"
)

func (h *WebappHandler) signUp(c *gin.Context) {
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
	err := h.services.SignupLocal(c.Request.Context(), req)
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

func (h *WebappHandler) signIn(c *gin.Context) {
	var req domain.SignInReq
	if err := c.ShouldBind(&req); err != nil {
		result := utils.GenerateMessage(err, reflect.TypeOf(req))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "invalid request",
			Errors:  result,
		})
		return
	}

	res, err := h.services.Signin(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}
	session := sessions.Default(c)
	session.Set(domain.SESSION_ID, res.UserWithIdentity.User.ID.String())
	session.Set(domain.SESSION_ROLE, res.UserWithIdentity.User.Role)

	if err := session.Save(); err != nil {
		h.log.Error(err, "failed to save session")
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrInternalServerError.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
		Data:    res,
	})
}

func (h *WebappHandler) me(c *gin.Context) {
	userID, err := utils.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	res, err := h.services.Me(c.Request.Context(), userID)
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
