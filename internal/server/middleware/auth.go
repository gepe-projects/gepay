package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ilhamgepe/gepay/internal/domain"
)

func (mw *Middleware) WithTokenAuth(c *gin.Context) {
	// check token di header
	bearerToken := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if bearerToken == nil || len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrUnauthorized.Error(),
			},
		})
		c.Abort()
		return
	}

	// validate token
	claims, err := mw.security.VerifyToken(bearerToken[1], mw.config.JWT.Secret)
	if err != nil {
		mw.log.Debugf("failed to verify token: %v", err)
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrUnauthorized.Error(),
			},
		})
		c.Abort()
		return
	}

	c.Set(domain.CtxUserIDKey, claims.RegisteredClaims.Subject)

	c.Next()
}

func (mw *Middleware) WithSessionAuth(c *gin.Context) {
	// check session
	session := sessions.Default(c)
	userID := session.Get(domain.SESSION_ID)
	Role := session.Get(domain.SESSION_ROLE)
	if userID == nil || Role == nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrUnauthorized.Error(),
			},
		})
		c.Abort()
		return
	}

	// memperpanjang lifetime session
	// if err := session.Save(); err != nil {
	// 	mw.log.Warnf("failed to extends session lifetime %v\n", err)
	// }

	mw.log.Infof("user: %s, role: %s", userID, Role)
	c.Set(domain.CtxUserIDKey, userID)
	c.Next()
}
