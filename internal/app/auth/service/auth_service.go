package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/internal/server/security"
	"github.com/ilhamgepe/gepay/pkg/config"
	"github.com/ilhamgepe/gepay/pkg/logger"
	uow "github.com/ilhamgepe/gepay/pkg/unitOfWork"
	"github.com/ilhamgepe/gepay/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type authService struct {
	userService     domain.UserService
	merchantService domain.MerchantService
	log             logger.Logger
	security        *security.Security
	db              *sqlx.DB
	config          config.App
}

var (
	readmeClaim = "Hire me, https://github.com/ilhamgepe"
	issuerClaim = "gepay"
)

func NewAuthService(userService domain.UserService, merchantService domain.MerchantService, log logger.Logger, security *security.Security, db *sqlx.DB, config config.App) domain.AuthService {
	return &authService{
		userService:     userService,
		merchantService: merchantService,
		log:             log,
		security:        security,
		db:              db,
		config:          config,
	}
}

func (s *authService) Signin(ctx context.Context, req domain.SignInReq) (res domain.SigninRes, err error) {
	// login limter
	delay, err := s.security.CheckBan(ctx, req.Email)
	if err != nil {
		return res, domain.ErrInternalServerError
	}

	if delay > 0 {
		return res, fmt.Errorf("too many attempts, please try again after %s", delay.String())
	}

	// check email
	user, err := s.userService.FindUserWithIdentityByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return res, domain.ErrInvalidCredentials
		}
		return res, err
	}

	// compare password
	if !utils.CheckPassword(req.Password, *user.UserIdentity.PasswordHash) {
		_, _ = s.security.IncrementAttempts(ctx, req.Email)
		return res, domain.ErrInvalidCredentials
	}
	_ = s.security.ResetAuthAttempts(ctx, req.Email)

	// generate token
	// AT
	claims := domain.TokenCLaims{
		Role:   user.User.Role,
		Readme: readmeClaim,
		Scopes: []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuerClaim,
			Subject:   user.User.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenExpiresIn)),
		},
	}

	accessToken, err := s.security.GenerateToken(claims, s.config.JWT.Secret)
	if err != nil {
		s.log.Error(err, "error generate token")
		return res, domain.ErrInternalServerError
	}

	RefreshClaims := domain.TokenCLaims{
		Role:   user.User.Role,
		Readme: readmeClaim,
		Scopes: []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuerClaim,
			Subject:   user.User.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshTokenExpiresIn)),
			ID:        uuid.New().String(),
		},
	}

	refreshToken, err := s.security.GenerateToken(RefreshClaims, s.config.JWT.RefreshSecret)
	if err != nil {
		s.log.Error(err, "error generate refresh token")
		return res, domain.ErrInternalServerError
	}

	return domain.SigninRes{
		Token:         accessToken,
		Refresh_token: refreshToken,
		User:          user.User,
	}, nil
}

func (s *authService) SignupLocal(ctx context.Context, req domain.SignUpReq) error {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hash

	uow := uow.NewUnitOfWork(s.db)
	err = uow.Do(ctx, func(tx *sqlx.Tx) error {
		// create user
		id, err := s.userService.CreateUserTx(ctx, tx, req)
		if err != nil {
			return err
		}

		// create user identity
		req.UserID = id
		err = s.userService.CreateUserIdentityTx(ctx, tx, req)
		if err != nil {
			return err
		}

		// create merchant
		merchantReq := domain.CreateMerchantReq{
			UserID:       req.UserID,
			OwnerName:    req.OwnerName,
			BusinessName: req.BusinessName,
		}
		err = s.merchantService.CreateMerchantTx(ctx, tx, merchantReq)

		return err
	})

	return err
}

func (s *authService) Me(ctx context.Context, id uuid.UUID) (res domain.UserWithIdentity, err error) {
	res, err = s.userService.FindUserWithIdentityByID(ctx, id)
	if err != nil {
		return
	}

	return
}

// Refresh implements domain.AuthService.
func (s *authService) Refresh(ctx context.Context, token string) (domain.RefreshRes, error) {
	claims, err := s.security.VerifyToken(token, s.config.JWT.RefreshSecret)
	if err != nil {
		return domain.RefreshRes{}, domain.ErrUnauthorized
	}

	// generate token
	// AT
	AccessClaim := domain.TokenCLaims{
		Role:   claims.Role,
		Readme: readmeClaim,
		Scopes: []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuerClaim,
			Subject:   claims.RegisteredClaims.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.AccessTokenExpiresIn)),
		},
	}
	accessToken, err := s.security.GenerateToken(AccessClaim, s.config.JWT.Secret)
	if err != nil {
		s.log.Error(err, "error generate token")
		return domain.RefreshRes{}, domain.ErrInternalServerError
	}

	refreshClaims := domain.TokenCLaims{
		Role:   claims.Role,
		Readme: readmeClaim,
		Scopes: []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuerClaim,
			Subject:   claims.RegisteredClaims.Subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.RefreshTokenExpiresIn)),
		},
	}

	refreshToken, err := s.security.GenerateToken(refreshClaims, s.config.JWT.RefreshSecret)
	if err != nil {
		s.log.Error(err, "error generate refresh token")
		return domain.RefreshRes{}, domain.ErrInternalServerError
	}
	return domain.RefreshRes{
		Token:         accessToken,
		Refresh_token: refreshToken,
	}, nil
}
