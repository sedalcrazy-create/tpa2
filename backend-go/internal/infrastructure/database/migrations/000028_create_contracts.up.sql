-- Create contract_types table
CREATE TABLE IF NOT EXISTS contract_types (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Code and naming
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,

    -- Category
    category VARCHAR(50),

    -- Default settings
    default_duration_months INT,
    default_renewal_type VARCHAR(20),
    default_grace_period_days INT,

    -- Template
    contract_template TEXT,
    terms_template TEXT,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_contract_type_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT uq_contract_type_code_tenant UNIQUE (code, tenant_id, deleted_at)
);

-- Create contracts table
CREATE TABLE IF NOT EXISTS contracts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Contract identification
    contract_number VARCHAR(50) NOT NULL,
    contract_type_id BIGINT NOT NULL,

    -- Parties - Insurer
    insurer_name VARCHAR(200) NOT NULL,
    insurer_code VARCHAR(50),
    insurer_rep_name VARCHAR(200),
    insurer_rep_title VARCHAR(100),

    -- Parties - Employer
    employer_name VARCHAR(200) NOT NULL,
    employer_code VARCHAR(50),
    employer_national_id VARCHAR(11),
    employer_economic_code VARCHAR(14),
    employer_rep_name VARCHAR(200),
    employer_rep_title VARCHAR(100),
    employer_phone VARCHAR(20),
    employer_email VARCHAR(100),
    employer_address TEXT,

    -- Contract period
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    sign_date TIMESTAMP,
    effective_date TIMESTAMP,

    -- Contract details
    total_insured INT DEFAULT 0,
    main_insured INT DEFAULT 0,
    dependents_allowed INT DEFAULT 0,
    max_dependents_per_main INT DEFAULT 4,

    -- Financial terms
    total_premium_amount BIGINT,
    premium_per_person BIGINT,
    payment_method VARCHAR(50),
    payment_day INT,
    advance_payment_percent DECIMAL(5,2),
    advance_payment_amount BIGINT,

    -- Coverage limits
    annual_coverage_limit BIGINT,
    per_claim_limit BIGINT,
    franchise_amount BIGINT,
    franchise_percentage DECIMAL(5,2),

    -- Renewal terms
    renewal_type VARCHAR(20) DEFAULT 'MANUAL',
    grace_period_days INT DEFAULT 30,
    renewal_notify_days INT DEFAULT 60,
    allow_early_termination BOOLEAN DEFAULT FALSE,
    termination_penalty BIGINT,

    -- Addendums
    addendum_count INT DEFAULT 0,
    last_addendum_date TIMESTAMP,
    last_addendum_no VARCHAR(50),

    -- Documents
    contract_file_path VARCHAR(500),
    signed_contract_path VARCHAR(500),
    addendums_path TEXT,
    attachments_path TEXT,

    -- Approval workflow
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    approved_by BIGINT,
    approved_at TIMESTAMP,
    approval_notes TEXT,

    -- Termination
    termination_date TIMESTAMP,
    termination_reason TEXT,
    terminated_by BIGINT,

    -- Notes
    terms TEXT,
    special_conditions TEXT,
    notes TEXT,
    internal_notes TEXT,

    -- System tracking
    created_by BIGINT,
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_contract_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_contract_type FOREIGN KEY (contract_type_id) REFERENCES contract_types(id) ON DELETE RESTRICT,
    CONSTRAINT fk_contract_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_contract_terminated_by FOREIGN KEY (terminated_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_contract_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT uq_contract_number_tenant UNIQUE (contract_number, tenant_id, deleted_at)
);

-- Create indexes for contract_types
CREATE INDEX idx_contract_type_tenant ON contract_types(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_type_category ON contract_types(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_type_is_active ON contract_types(is_active) WHERE deleted_at IS NULL;

-- Create indexes for contracts
CREATE INDEX idx_contract_tenant ON contracts(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_type ON contracts(contract_type_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_employer_name ON contracts(employer_name) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_employer_code ON contracts(employer_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_employer_national ON contracts(employer_national_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_employer_economic ON contracts(employer_economic_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_start_date ON contracts(start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_end_date ON contracts(end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_sign_date ON contracts(sign_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_effective_date ON contracts(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_status ON contracts(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_contract_is_active ON contracts(is_active) WHERE deleted_at IS NULL;

-- Add comments
COMMENT ON TABLE contract_types IS 'Types of insurance contracts (individual, group, corporate, etc.)';
COMMENT ON TABLE contracts IS 'Insurance contracts between insurer and employer/individual';
