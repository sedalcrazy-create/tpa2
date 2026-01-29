-- Create provider_infos table
CREATE TABLE IF NOT EXISTS provider_infos (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Provider identification
    national_code VARCHAR(10) NOT NULL UNIQUE,
    medical_code VARCHAR(20) UNIQUE,
    pharmacy_code VARCHAR(20) UNIQUE,
    license_number VARCHAR(50),
    insurance_code VARCHAR(20),

    -- Personal information
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    father_name VARCHAR(100),
    birth_date VARCHAR(10),
    gender VARCHAR(10),

    -- Professional information
    provider_type VARCHAR(50) NOT NULL,
    specialty VARCHAR(100),
    sub_specialty VARCHAR(100),
    academic_degree VARCHAR(50),

    -- Education
    university_name VARCHAR(200),
    graduation_year INT,
    medical_council_name VARCHAR(100),

    -- License and certification
    license_issue_date VARCHAR(10),
    license_expiry_date VARCHAR(10),
    certification_body VARCHAR(200),
    is_license_valid BOOLEAN DEFAULT TRUE,

    -- Contact information
    phone VARCHAR(20),
    mobile VARCHAR(20),
    email VARCHAR(100),
    website VARCHAR(200),

    -- Address
    province_id BIGINT,
    city_id BIGINT,
    address TEXT,
    postal_code VARCHAR(10),

    -- Profile
    biography TEXT,
    photo_path VARCHAR(500),
    signature_path VARCHAR(500),
    achievements TEXT,
    research_interests TEXT,
    languages VARCHAR(200),

    -- Statistics
    total_prescriptions INT DEFAULT 0,
    total_claims INT DEFAULT 0,
    avg_claim_amount BIGINT,
    last_prescription_date VARCHAR(10),

    -- Verification
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at VARCHAR(10),
    verified_by BIGINT,
    verification_notes TEXT,

    -- Flags and status
    is_active BOOLEAN DEFAULT TRUE,
    is_blacklisted BOOLEAN DEFAULT FALSE,
    blacklist_reason TEXT,
    blacklisted_at VARCHAR(10),
    is_suspended BOOLEAN DEFAULT FALSE,
    suspension_reason TEXT,
    suspended_at VARCHAR(10),

    -- Notes
    notes TEXT,
    internal_notes TEXT,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_provider_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_provider_province FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE SET NULL,
    CONSTRAINT fk_provider_city FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE SET NULL,
    CONSTRAINT fk_provider_verified_by FOREIGN KEY (verified_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create provider_center_mappings table
CREATE TABLE IF NOT EXISTS provider_center_mappings (
    id BIGSERIAL PRIMARY KEY,

    provider_id BIGINT NOT NULL,
    center_id BIGINT NOT NULL,

    -- Mapping properties
    role VARCHAR(50),
    start_date VARCHAR(10),
    end_date VARCHAR(10),
    is_active BOOLEAN DEFAULT TRUE,
    is_primary BOOLEAN DEFAULT FALSE,
    work_schedule TEXT,
    notes TEXT,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_pcm_provider FOREIGN KEY (provider_id) REFERENCES provider_infos(id) ON DELETE CASCADE,
    CONSTRAINT fk_pcm_center FOREIGN KEY (center_id) REFERENCES centers(id) ON DELETE CASCADE,
    CONSTRAINT uq_pcm_provider_center UNIQUE (provider_id, center_id, deleted_at)
);

-- Create indexes for provider_infos
CREATE INDEX idx_provider_tenant ON provider_infos(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_national ON provider_infos(national_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_medical ON provider_infos(medical_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_pharmacy ON provider_infos(pharmacy_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_license ON provider_infos(license_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_insurance_code ON provider_infos(insurance_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_gender ON provider_infos(gender) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_type ON provider_infos(provider_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_specialty ON provider_infos(specialty) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_is_license_valid ON provider_infos(is_license_valid) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_mobile ON provider_infos(mobile) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_email ON provider_infos(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_province ON provider_infos(province_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_city ON provider_infos(city_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_is_verified ON provider_infos(is_verified) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_is_active ON provider_infos(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_is_blacklisted ON provider_infos(is_blacklisted) WHERE deleted_at IS NULL;
CREATE INDEX idx_provider_is_suspended ON provider_infos(is_suspended) WHERE deleted_at IS NULL;

-- Create indexes for provider_center_mappings
CREATE INDEX idx_pcm_provider ON provider_center_mappings(provider_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pcm_center ON provider_center_mappings(center_id) WHERE deleted_at IS NULL;

-- Add comments
COMMENT ON TABLE provider_infos IS 'Healthcare providers information (physicians, pharmacists, etc.)';
COMMENT ON TABLE provider_center_mappings IS 'Many-to-many mapping between providers and centers';
