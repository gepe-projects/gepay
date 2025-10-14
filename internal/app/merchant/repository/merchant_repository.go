package repository

import (
	"context"

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

// CreateMerchantTx implements domain.MerchantRepository.
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
