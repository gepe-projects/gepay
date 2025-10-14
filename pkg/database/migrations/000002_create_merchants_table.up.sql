CREATE TABLE merchants (
  id UUID DEFAULT uuidv7() PRIMARY KEY,
  user_id UUID NOT NULL,
  owner_name VARCHAR(150) NOT NULL,
  business_name VARCHAR(150) NOT NULL, -- nama brand/toko
  business_type VARCHAR(100),
  description TEXT,
  logo_url VARCHAR(255),

  -- Status
  status VARCHAR(50) NOT NULL DEFAULT 'pending', -- 'pending', 'active', 'suspended'
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  disbursement_status BOOLEAN NOT NULL DEFAULT FALSE,

  country VARCHAR(100) DEFAULT 'ID',
  currency VARCHAR(10) DEFAULT 'IDR',

  -- Setting internal platform
  webhook_url VARCHAR(255),
  webhook_secret VARCHAR(100),
  is_test_mode BOOLEAN NOT NULL DEFAULT TRUE,

  metadata JSONB DEFAULT '{}'::jsonb,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_merchants_user_id ON merchants(user_id);
CREATE INDEX idx_merchants_status ON merchants(status);
CREATE INDEX idx_merchant_owner_name ON merchants(owner_name);
CREATE INDEX idx_merchant_Business_name ON merchants(business_name);
