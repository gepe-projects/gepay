package service

import (
	"context"

	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type merchantService struct {
	merchantRepo domain.MerchantRepository
	log          logger.Logger
}

func NewMerchantService(merchantRepo domain.MerchantRepository, log logger.Logger) domain.MerchantService {
	return &merchantService{
		merchantRepo: merchantRepo,
		log:          log,
	}
}

// CreateMerchantTx implements domain.MerchantService.
func (m *merchantService) CreateMerchantTx(ctx context.Context, tx *sqlx.Tx, req domain.CreateMerchantReq) error {
	return m.merchantRepo.CreateMerchantTx(ctx, tx, req)
}
