package domain

import (
	"time"

	"github.com/google/uuid"
)

// ! TODO : implement file upload
type MerchantKYCReq struct {
	UserID       uuid.UUID
	MerchantID   uuid.UUID
	KYCType      string `json:"kyc_type" form:"kyc_type" binding:"required,oneof=personal business"`
	IDCardNumber string `json:"id_card_number" form:"id_card_number" binding:"required,numeric,len=16"`
	TaxIDNumber  string `json:"tax_id_number" form:"tax_id_number" binding:"required,numeric,len=16"`
	// TaxIDImage   *multipart.FileHeader `json:"tax_id_image" form:"tax_id_image"`

	// Legalitas Usaha (business only)
	LegalBusinessName     string
	BusinessLicenseNumber string `json:"business_license_number" form:"business_license_number"`
	// BusinessLicenseImage    *multipart.FileHeader `json:"business_license_image" form:"business_license_image" binding:"omitempty,required_if=KYCType business"`
	DeedNumber string `json:"deed_number" form:"deed_number"`
	// DeedImage  *multipart.FileHeader `json:"deed_image" form:"deed_image" binding:"omitempty,required_if=KYCType business"`
	Address    string `json:"address" form:"address" binding:"required,max=255,min=8"`
	WebsiteURL string `json:"website_url" form:"website_url"`
}

type MerchantKYC struct {
	// BASIC & IDENTIFICATION
	MerchantID uuid.UUID `json:"merchant_id" db:"merchant_id"`
	KYCType    string    `json:"kyc_type" db:"kyc_type"` // 'personal' atau 'business'

	// PERSONAL AND BUSINESS
	IDCardNumber string `json:"id_card_number" db:"id_card_number"` // CHAR(16)
	IDCardURL    string `json:"id_card_url" db:"id_card_url"`       // TEXT NOT NULL

	// NPWP
	TaxIDNumber *string `json:"tax_id_number" db:"tax_id_number"` // CHAR(15) (NPWP)
	TaxIDURL    *string `json:"tax_id_url" db:"tax_id_url"`

	// BUSINESS ONLY (NULLABLE)
	LegalBusinessName     *string `json:"legal_business_name" db:"legal_business_name"`
	BusinessLicenseNumber *string `json:"business_license_number" db:"business_license_number"` // VARCHAR(20) (NIB/ijin usaha)
	BusinessLicenseURL    *string `json:"business_license_url" db:"business_license_url"`

	DeedNumber *string `json:"deed_number" db:"deed_number"`
	DeedURL    *string `json:"deed_url" db:"deed_url"`
	Address    *string `json:"address" db:"address"`
	WebsiteURL *string `json:"website_url" db:"website_url"`

	// STATUS & METADATA
	Status          string        `json:"status" db:"status"` // 'pending', 'in_review', 'verified', 'rejected'
	RejectionReason *string       `json:"rejection_reason" db:"rejection_reason"`
	VerifiedAt      *time.Time    `json:"verified_at" db:"verified_at"`
	ReviewedBy      uuid.NullUUID `json:"reviewed_by" db:"reviewed_by"`

	// TIMESTAMPS
	Metadata  []byte    `json:"metadata" db:"metadata"` // JSONB di PostgreSQL
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
