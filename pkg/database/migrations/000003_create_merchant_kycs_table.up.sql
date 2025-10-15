CREATE TABLE merchant_kycs (
    merchant_id UUID PRIMARY KEY,
    
    -- KYC TYPE
    kyc_type VARCHAR(20) NOT NULL CHECK (kyc_type IN ('personal', 'business')),
    -- Penambahan CHECK Constraint untuk memastikan nilai hanya 'personal' atau 'business'
    
    -- PERSONAL AND BUSINESS ALSO
    -- NIK/KTP umumnya 16 digit. Dibuat CHAR(16) untuk konsisten.
    id_card_number CHAR(16) NOT NULL,
    -- Disarankan menggunakan TEXT untuk URL dokumen yang diunggah
    id_card_url TEXT NOT NULL, -- KTP wajib untuk semua
    
    -- NPWP (15-16 digit)
    -- Dibuat CHAR(20) untuk jaga jaga aja wkwkw.
    tax_id_number CHAR(20), 
    tax_id_url TEXT,

    -- BUSINESS ONLY (Dibuat NULLABLE agar bisa digunakan oleh 'personal' tanpa error)
    legal_business_name VARCHAR(150),
    
    -- Nomor Izin Usaha/NIB (sering 13/16 digit)
    business_license_number VARCHAR(20),
    business_license_url TEXT,
    
    -- Akta Pendirian
    deed_number VARCHAR(100),
    deed_url TEXT,
    
    address TEXT, -- Alamat bisa sangat panjang
    website_url VARCHAR(255),

    -- STATUS AND METADATA
    -- Menggunakan ENUM atau constraint CHECK lebih baik daripada VARCHAR untuk status
    status VARCHAR(20) NOT NULL DEFAULT 'pending' 
        CHECK (status IN ('pending', 'in_review', 'verified', 'rejected')),
        
    rejection_reason TEXT,
    verified_at TIMESTAMP,
    reviewed_by UUID,

    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    
    -- INDEXING
    -- Membuat index pada kolom yang sering digunakan untuk filter atau lookup.
    
    -- FOREIGN KEY
    FOREIGN KEY(merchant_id) REFERENCES merchants(id) ON DELETE CASCADE
);

-- Indexing untuk mempercepat query status
CREATE INDEX idx_merchant_kycs_status ON merchant_kycs (status);