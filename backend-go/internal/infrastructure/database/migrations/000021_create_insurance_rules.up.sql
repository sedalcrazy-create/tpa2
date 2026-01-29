-- Create insurance_rules table
CREATE TABLE IF NOT EXISTS insurance_rules (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    insurance_id BIGINT NOT NULL,

    -- Rule identification
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    rule_type VARCHAR(50) NOT NULL,

    -- Coverage limits
    annual_limit BIGINT,
    per_claim_limit BIGINT,
    lifetime_limit BIGINT,
    daily_limit BIGINT,
    monthly_limit BIGINT,
    per_service_limit BIGINT,
    per_drug_limit BIGINT,
    hospitalization_days INT,

    -- Service-specific limits
    drug_limit BIGINT,
    dental_limit BIGINT,
    optical_limit BIGINT,
    physiotherapy_limit BIGINT,
    lab_limit BIGINT,
    imaging_limit BIGINT,
    surgery_limit BIGINT,

    -- Deductible
    deductible_amount BIGINT,
    deductible_percentage DECIMAL(5,2),
    deductible_type VARCHAR(20),

    -- Waiting periods
    general_waiting_days INT,
    dental_waiting_days INT,
    optical_waiting_days INT,
    maternity_waiting_days INT,
    surgery_waiting_days INT,

    -- Co-payment
    co_payment_percentage DECIMAL(5,2),
    co_payment_amount BIGINT,

    -- Age/Gender restrictions
    min_age INT,
    max_age INT,
    gender_limit VARCHAR(10),
    appliesto_main BOOLEAN DEFAULT TRUE,
    applies_to_dep BOOLEAN DEFAULT TRUE,
    dependent_limit INT,
    child_age_limit INT,

    -- Special conditions
    requires_pre_auth BOOLEAN DEFAULT FALSE,
    requires_prescription BOOLEAN DEFAULT FALSE,
    requires_referal BOOLEAN DEFAULT FALSE,
    allow_emergency BOOLEAN DEFAULT TRUE,
    allow_chronic_illness BOOLEAN DEFAULT TRUE,
    allow_pre_existing BOOLEAN DEFAULT FALSE,
    network_only BOOLEAN DEFAULT FALSE,

    -- Exclusions
    excluded_services TEXT,
    excluded_drugs TEXT,
    excluded_icd10 TEXT,

    -- Renewal rules
    auto_renewal BOOLEAN DEFAULT FALSE,
    renewal_grace_days INT,
    premium_increase DECIMAL(5,2),

    -- Priority and status
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    -- Validity period
    effective_date TIMESTAMP,
    expiry_date TIMESTAMP,

    -- Rule engine integration
    rule_engine_code TEXT,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_insurance_rule_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_insurance_rule_insurance FOREIGN KEY (insurance_id) REFERENCES insurances(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_insurance_rule_tenant ON insurance_rules(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_insurance ON insurance_rules(insurance_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_code ON insurance_rules(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_type ON insurance_rules(rule_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_priority ON insurance_rules(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_is_active ON insurance_rules(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_effective ON insurance_rules(effective_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_insurance_rule_expiry ON insurance_rules(expiry_date) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE insurance_rules IS 'Business rules for insurance policies - coverage limits, deductibles, waiting periods';
