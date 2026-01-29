-- Create item_price_conditions table
CREATE TABLE IF NOT EXISTS item_price_conditions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Item reference
    item_id BIGINT,
    category_id BIGINT,
    sub_category_id BIGINT,
    group_id BIGINT,

    -- Insurance reference
    insurance_id BIGINT,

    -- Pricing rules
    coverage_percentage DECIMAL(5,2),
    max_coverage_amount BIGINT,
    min_coverage_amount BIGINT,
    fixed_amount BIGINT,

    -- Franchise
    franchise_percentage DECIMAL(5,2),
    franchise_amount BIGINT,
    max_franchise BIGINT,

    -- Deductible
    base_insurance_share DECIMAL(5,2),

    -- Quantity limits
    max_quantity_per_day INT,
    max_quantity_per_month INT,
    max_quantity_per_year INT,

    -- Age/Gender filters
    min_age INT,
    max_age INT,
    gender VARCHAR(10),
    is_for_main BOOLEAN DEFAULT TRUE,
    is_for_dep BOOLEAN DEFAULT TRUE,

    -- Waiting period
    waiting_period_days INT,

    -- Special conditions
    needs_prescription BOOLEAN DEFAULT FALSE,
    needs_pre_approval BOOLEAN DEFAULT FALSE,
    needs_medical_opinion BOOLEAN DEFAULT FALSE,
    requires_diagnosis BOOLEAN DEFAULT FALSE,
    allowed_for_chronic_ill BOOLEAN DEFAULT TRUE,

    -- Priority and status
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    -- Validity period
    start_date TIMESTAMP,
    end_date TIMESTAMP,

    -- Description
    title VARCHAR(200),
    description TEXT,
    note TEXT,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_ipc_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_ipc_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE CASCADE,
    CONSTRAINT fk_ipc_category FOREIGN KEY (category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_ipc_sub_category FOREIGN KEY (sub_category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_ipc_group FOREIGN KEY (group_id) REFERENCES item_groups(id) ON DELETE SET NULL,
    CONSTRAINT fk_ipc_insurance FOREIGN KEY (insurance_id) REFERENCES insurances(id) ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_ipc_tenant ON item_price_conditions(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_item ON item_price_conditions(item_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_category ON item_price_conditions(category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_subcategory ON item_price_conditions(sub_category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_group ON item_price_conditions(group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_insurance ON item_price_conditions(insurance_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_priority ON item_price_conditions(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_is_active ON item_price_conditions(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_start_date ON item_price_conditions(start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_ipc_end_date ON item_price_conditions(end_date) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE item_price_conditions IS 'Pricing rules and conditions for drugs/services - Core pricing engine';
