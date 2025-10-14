package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type userService struct {
	log      logger.Logger
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository, log logger.Logger) domain.UserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
	}
}

// CreateUserIdentityTx implements domain.UserService.
func (u *userService) CreateUserIdentityTx(ctx context.Context, tx *sqlx.Tx, req domain.SignUpReq) error {
	return u.userRepo.CreateUserIdentityTx(ctx, tx, req)
}

// CreateUserTx implements domain.UserService.
func (u *userService) CreateUserTx(ctx context.Context, tx *sqlx.Tx, req domain.SignUpReq) (uuid.UUID, error) {
	return u.userRepo.CreateUserTx(ctx, tx, req)
}

func (u *userService) FindUserWithIdentityByEmail(ctx context.Context, email string) (domain.UserWithIdentity, error) {
	return u.userRepo.FindUserWithIdentityByEmail(ctx, email)
}

func (u *userService) FindUserWithIdentityByID(ctx context.Context, id uuid.UUID) (domain.UserWithIdentity, error) {
	return u.userRepo.FindUserWithIdentityByID(ctx, id)
}
