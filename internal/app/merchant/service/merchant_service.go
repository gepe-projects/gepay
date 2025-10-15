package service

import (
	"context"

	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
	uow "github.com/ilhamgepe/gepay/pkg/unitOfWork"
	"github.com/jmoiron/sqlx"
)

type merchantService struct {
	merchantRepo domain.MerchantRepository
	log          logger.Logger
	db           *sqlx.DB
}

func NewMerchantService(merchantRepo domain.MerchantRepository, db *sqlx.DB, log logger.Logger) domain.MerchantService {
	return &merchantService{
		merchantRepo: merchantRepo,
		log:          log,
		db:           db,
	}
}

// CreateMerchantKYC implements domain.MerchantService.
func (m *merchantService) CreateMerchantKYC(ctx context.Context, req domain.MerchantKYCReq) error {
	merchant, err := m.merchantRepo.GetMerchantByUserID(ctx, req.UserID)
	if err != nil {
		return err
	}

	uow := uow.NewUnitOfWork(m.db)

	err = uow.Do(ctx, func(tx *sqlx.Tx) error {
		// create merchant kyc
		req.MerchantID = merchant.ID
		req.LegalBusinessName = merchant.BusinessName
		err := m.merchantRepo.CreateMerchantKYCTx(ctx, tx, req)
		if err != nil {
			return err
		}
		updateMerchantReq := domain.UpdateMerchantReq{
			MerchantID:         merchant.ID,
			OwnerName:          merchant.OwnerName,
			BusinessName:       req.LegalBusinessName,
			BusinessType:       &req.KYCType,
			Description:        nil,
			Status:             merchant.Status,
			Verified:           merchant.Verified,
			DisbursementStatus: merchant.DisbursementStatus,
			Country:            merchant.Country,
			Currency:           merchant.Currency,
			WebhookURL:         merchant.WebhookURL,
			WebhookSecret:      merchant.WebhookSecret,
			IsTestMode:         merchant.IsTestMode,
			Metadata:           merchant.Metadata,
		}
		err = m.merchantRepo.UpdateMerchantTx(ctx, tx, updateMerchantReq)

		return err
	})

	return err
}

// CreateMerchantTx implements domain.MerchantService.
func (m *merchantService) CreateMerchantTx(ctx context.Context, tx *sqlx.Tx, req domain.CreateMerchantReq) error {
	return m.merchantRepo.CreateMerchantTx(ctx, tx, req)
}
