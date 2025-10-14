CREATE TABLE merchant_kyc (
  merchant_id UUID PRIMARY KEY,
  kyc_type VARCHAR(50) NOT NULL, -- 'personal' | 'business'

  -- PERSONAL AND BUSINESS ALSO
  id_card_number VARCHAR(100) NOT NULL,
  id_card_url VARCHAR(255),
  npwp_number VARCHAR(100),
  npwp_url VARCHAR(255),

  -- BUSINESS
  legal_business_name VARCHAR(150),
  business_license_no VARCHAR(100),
  business_license_url VARCHAR(255),
  tax_id_number VARCHAR(100),
  tax_id_url VARCHAR(255),
  deed_number VARCHAR(100),
  deed_url VARCHAR(255),
  address TEXT,
  website_url VARCHAR(255),

  -- STATUS
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  rejection_reason TEXT,
  verified_at TIMESTAMP,
  reviewed_by UUID,

  metadata JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),

  FOREIGN KEY(merchant_id) REFERENCES merchants(id) ON DELETE CASCADE
);
