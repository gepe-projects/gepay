package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
)

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID, ok := c.Get(domain.CtxUserIDKey)
	if !ok {
		return uuid.UUID{}, domain.ErrInternalServerError
	}

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		return uuid.UUID{}, domain.ErrInternalServerError
	}
	return uid, nil
}
