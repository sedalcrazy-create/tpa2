-- Create condition_groups table
CREATE TABLE IF NOT EXISTS condition_groups (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Group information
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,

    -- Logical operator
    logic_operator VARCHAR(10) DEFAULT 'AND',

    -- Priority
    priority INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_cond_group_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT uq_cond_group_code_tenant UNIQUE (code, tenant_id, deleted_at)
);

-- Create condition_group_mappings table
CREATE TABLE IF NOT EXISTS condition_group_mappings (
    id BIGSERIAL PRIMARY KEY,

    condition_group_id BIGINT NOT NULL,
    item_price_condition_id BIGINT NOT NULL,

    -- Mapping properties
    sort_order INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_cgm_group FOREIGN KEY (condition_group_id) REFERENCES condition_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_cgm_condition FOREIGN KEY (item_price_condition_id) REFERENCES item_price_conditions(id) ON DELETE CASCADE,
    CONSTRAINT uq_cgm_group_condition UNIQUE (condition_group_id, item_price_condition_id, deleted_at)
);

-- Create indexes
CREATE INDEX idx_cond_group_tenant ON condition_groups(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cond_group_priority ON condition_groups(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_cond_group_is_active ON condition_groups(is_active) WHERE deleted_at IS NULL;

CREATE INDEX idx_cgm_group ON condition_group_mappings(condition_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_cgm_condition ON condition_group_mappings(item_price_condition_id) WHERE deleted_at IS NULL;

-- Add comments
COMMENT ON TABLE condition_groups IS 'Groups of pricing conditions with logical operators (AND/OR)';
COMMENT ON TABLE condition_group_mappings IS 'Many-to-many mapping between condition groups and price conditions';
