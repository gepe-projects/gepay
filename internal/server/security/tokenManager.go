package security

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamgepe/gepay/internal/domain"
)

func (s *Security) GenerateToken(claims domain.TokenCLaims, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", domain.ErrUnauthorized
	}
	return tokenString, err
}

func (s *Security) VerifyToken(token string, secret string) (domain.TokenCLaims, error) {
	var claim domain.TokenCLaims
	parsedToken, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		s.log.Error(err, "failed to parse token")
		return claim, domain.ErrUnauthorized
	}

	if !parsedToken.Valid {
		s.log.Error(err, "invalid token")
		return claim, domain.ErrUnauthorized
	}

	return claim, nil
}
