package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/ilhamgepe/gepay/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type merchantRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewMerchantRepo(db *sqlx.DB, log logger.Logger) domain.MerchantRepository {
	return &merchantRepository{
		db:  db,
		log: log,
	}
}

func (m *merchantRepository) CreateMerchantTx(ctx context.Context, tx *sqlx.Tx, req domain.CreateMerchantReq) error {
	query := `
		INSERT INTO merchants (user_id,owner_name,business_name)
		VALUES ($1, $2, $3)
	`

	_, err := tx.ExecContext(ctx, query, req.UserID, req.OwnerName, req.BusinessName)
	if err != nil {
		m.log.Errorf(err, "failed to create merchant")
		if utils.IsForeignKeyViolation(err) {
			return domain.ErrMerchantAlreadyExists
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (m *merchantRepository) GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (domain.Merchant, error) {
	var merchant domain.Merchant
	query := `
		SELECT
			id,user_id,owner_name,
			business_name,business_type,
			description,logo_url,status,
			verified,disbursement_status,country,
			currency,webhook_url,
			webhook_secret,is_test_mode,
			metadata
		FROM merchants
		WHERE user_id = $1
	`
	err := m.db.GetContext(ctx, &merchant, query, userID)
	if err != nil {
		m.log.Errorf(err, "failed to get merchant by user id")
		return domain.Merchant{}, domain.ErrInternalServerError
	}

	return merchant, nil
}

func (m *merchantRepository) UpdateMerchantTx(ctx context.Context, tx *sqlx.Tx, req domain.UpdateMerchantReq) error {
	query := `
		UPDATE merchants
		SET
			owner_name = :owner_name,
			business_name = :business_name,
			business_type = :business_type,
			description = :description,
			status = :status,
			verified = :verified,
			disbursement_status = :disbursement_status,
			country = :country,
			currency = :currency,
			webhook_url = :webhook_url,
			webhook_secret = :webhook_secret,
			is_test_mode = :is_test_mode,
			metadata = :metadata
		WHERE id = :merchant_id
	`

	_, err := tx.NamedExecContext(ctx, query, req)
	if err != nil {
		m.log.Errorf(err, "failed to update merchant")
		return domain.ErrInternalServerError
	}

	return nil
}

// MERCHANT KYCS
func (m *merchantRepository) CreateMerchantKYCTx(ctx context.Context, tx *sqlx.Tx, req domain.MerchantKYCReq) error {
	var query string
	var args []any

	args = append(args, req.MerchantID, req.KYCType, req.IDCardNumber, "https://ui-avatars.com/api/?name="+req.LegalBusinessName)

	if req.KYCType == "personal" {
		query = `
            INSERT INTO merchant_kycs (merchant_id, kyc_type, id_card_number,id_card_url, address)
            VALUES ($1, $2, $3, $4, $5)
        `
		args = append(args, req.Address)
	} else {
		query = `
            INSERT INTO merchant_kycs 
                (merchant_id, kyc_type, id_card_number,id_card_url, tax_id_number, legal_business_name, business_license_number, deed_number, address, website_url)
            VALUES 
                ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10)
        `
		args = append(args, req.TaxIDNumber, req.LegalBusinessName, req.BusinessLicenseNumber,
			req.DeedNumber, req.Address, req.WebsiteURL)
	}

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		if utils.IsDuplicateUniqueViolation(err) {
			return domain.ErrMerchantKYCAlreadyExists
		}
		m.log.Errorf(err, "failed to create merchant kyc")
		return domain.ErrInternalServerError
	}

	return nil
}
