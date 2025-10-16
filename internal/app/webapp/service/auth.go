package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	uow "github.com/ilhamgepe/gepay/pkg/unitOfWork"
	"github.com/ilhamgepe/gepay/pkg/utils"
	"github.com/jmoiron/sqlx"
)

func (s *webappService) Signin(ctx context.Context, req domain.SignInReq) (res domain.SigninRes, err error) {
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

	// generate session data

	return domain.SigninRes{
		UserWithIdentity: domain.UserWithIdentity{
			User:         user.User,
			UserIdentity: user.UserIdentity,
		},
	}, nil
}

func (s *webappService) SignupLocal(ctx context.Context, req domain.SignUpReq) error {
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

func (s *webappService) Me(ctx context.Context, id uuid.UUID) (res domain.UserWithIdentity, err error) {
	res, err = s.userService.FindUserWithIdentityByID(ctx, id)
	if err != nil {
		return
	}

	return
}
