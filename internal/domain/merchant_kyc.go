package domain

import (
	"time"

	"github.com/google/uuid"
)

type MerchantKYC struct {
	MerchantID uuid.UUID `db:"merchant_id" json:"merchant_id"`
	KYCType    string    `db:"kyc_type" json:"kyc_type"` // 'personal' | 'business'

	// Dokumen Identitas
	IDCardNumber string  `db:"id_card_number" json:"id_card_number"`
	IDCardURL    *string `db:"id_card_url" json:"id_card_url,omitempty"`
	NPWPNumber   *string `db:"npwp_number" json:"npwp_number,omitempty"`
	NPWPURL      *string `db:"npwp_url" json:"npwp_url,omitempty"`

	// Legalitas Usaha (business only)
	LegalBusinessName  *string `db:"legal_business_name" json:"legal_business_name,omitempty"`
	BusinessLicenseNo  *string `db:"business_license_no" json:"business_license_no,omitempty"`
	BusinessLicenseURL *string `db:"business_license_url" json:"business_license_url,omitempty"`
	TaxIDNumber        *string `db:"tax_id_number" json:"tax_id_number,omitempty"`
	TaxIDURL           *string `db:"tax_id_url" json:"tax_id_url,omitempty"`
	DeedNumber         *string `db:"deed_number" json:"deed_number,omitempty"`
	DeedURL            *string `db:"deed_url" json:"deed_url,omitempty"`
	Address            *string `db:"address" json:"address,omitempty"`
	WebsiteURL         *string `db:"website_url" json:"website_url,omitempty"`

	// Status Verifikasi
	Status          string     `db:"status" json:"status"`
	RejectionReason *string    `db:"rejection_reason" json:"rejection_reason,omitempty"`
	VerifiedAt      *time.Time `db:"verified_at" json:"verified_at,omitempty"`
	ReviewedBy      *uuid.UUID `db:"reviewed_by" json:"reviewed_by,omitempty"`

	Metadata  JSONB     `db:"metadata" json:"metadata"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
