package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MerchantRepository interface {
	CreateMerchantTx(ctx context.Context, tx *sqlx.Tx, req CreateMerchantReq) error
	GetMerchantByUserID(ctx context.Context, userID uuid.UUID) (Merchant, error)
	UpdateMerchantTx(ctx context.Context, tx *sqlx.Tx, req UpdateMerchantReq) error
	// KYC
	CreateMerchantKYCTx(ctx context.Context, tx *sqlx.Tx, req MerchantKYCReq) error
}

type MerchantService interface {
	CreateMerchantTx(ctx context.Context, tx *sqlx.Tx, req CreateMerchantReq) error
	CreateMerchantKYC(ctx context.Context, req MerchantKYCReq) error
}

type CreateMerchantReq struct {
	UserID       uuid.UUID `json:"user_id" binding:"required,uuid"`
	OwnerName    string    `json:"owner_name" binding:"required"`
	BusinessName string    `json:"business_name" binding:"required"`
}

// TODO: implement file upload
type UpdateMerchantReq struct {
	MerchantID   uuid.UUID `db:"merchant_id"`
	OwnerName    string    `db:"owner_name" json:"owner_name" form:"owner_name" binding:"required"`
	BusinessName string    `db:"business_name" json:"business_name" form:"business_name" binding:"required"`
	BusinessType *string   `db:"business_type" json:"business_type" form:"business_type" binding:"required,oneof=personal business"`
	Description  *string   `db:"description" json:"description" form:"description" binding:"omitempty,max=255,min=8"`

	Status             string `db:"status"`
	Verified           bool   `db:"verified"`
	DisbursementStatus bool   `db:"disbursement_status" json:"disbursement_status"`

	Country  string `db:"country" json:"country"`
	Currency string `db:"currency" json:"currency"`

	WebhookURL    *string `db:"webhook_url" json:"webhook_url,omitempty"`
	WebhookSecret *string `db:"webhook_secret" json:"webhook_secret,omitempty"`
	IsTestMode    bool    `db:"is_test_mode" json:"is_test_mode"`

	Metadata JSONB `db:"metadata" json:"metadata"` // JSONB
}

type Merchant struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserID       uuid.UUID `db:"user_id" json:"user_id"`
	OwnerName    string    `db:"owner_name" json:"owner_name"`
	BusinessName string    `db:"business_name" json:"business_name"`
	BusinessType *string   `db:"business_type" json:"business_type,omitempty"`
	Description  *string   `db:"description" json:"description,omitempty"`
	LogoURL      *string   `db:"logo_url" json:"logo_url,omitempty"`

	Status             string `db:"status" json:"status"`
	Verified           bool   `db:"verified" json:"verified"`
	DisbursementStatus bool   `db:"disbursement_status" json:"disbursement_status"`

	Country  string `db:"country" json:"country"`
	Currency string `db:"currency" json:"currency"`

	WebhookURL    *string `db:"webhook_url" json:"webhook_url,omitempty"`
	WebhookSecret *string `db:"webhook_secret" json:"webhook_secret,omitempty"`
	IsTestMode    bool    `db:"is_test_mode" json:"is_test_mode"`

	Metadata  JSONB     `db:"metadata" json:"metadata"` // JSONB
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
