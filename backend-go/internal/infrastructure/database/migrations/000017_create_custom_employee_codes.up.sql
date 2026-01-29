-- Create custom_employee_codes table
CREATE TABLE IF NOT EXISTS custom_employee_codes (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Code information
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,

    -- Discount settings
    discount_percentage DECIMAL(5,2),
    discount_amount BIGINT,

    -- Price limit settings
    max_price_percentage DECIMAL(5,2),

    -- Special flags
    is_retired BOOLEAN DEFAULT FALSE,
    no_limitation BOOLEAN DEFAULT FALSE,
    special_group BOOLEAN DEFAULT FALSE,
    priority_access BOOLEAN DEFAULT FALSE,

    -- Validity period
    start_date TIMESTAMP,
    end_date TIMESTAMP,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_cec_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT uq_cec_code_tenant UNIQUE (code, tenant_id, deleted_at)
);

-- Create indexes
CREATE INDEX idx_cec_tenant ON custom_employee_codes(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cec_is_retired ON custom_employee_codes(is_retired) WHERE deleted_at IS NULL;
CREATE INDEX idx_cec_start_date ON custom_employee_codes(start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_cec_end_date ON custom_employee_codes(end_date) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE custom_employee_codes IS 'Special employee codes with custom pricing/discount rules (retired, VIP, etc.)';
