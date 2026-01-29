-- Create employee_special_discounts table
CREATE TABLE IF NOT EXISTS employee_special_discounts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,

    -- Discount information
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    discount_type VARCHAR(50) NOT NULL,

    -- Discount values
    discount_percentage DECIMAL(5,2),
    discount_amount BIGINT,
    max_discount_amount BIGINT,

    -- Scope filters
    apply_to_item_id BIGINT,
    apply_to_category_id BIGINT,
    apply_to_group_id BIGINT,
    apply_to_service_type VARCHAR(50),

    -- Application rules
    apply_on_base_price BOOLEAN DEFAULT TRUE,
    apply_after_insurance BOOLEAN DEFAULT FALSE,
    combine_with_others BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 0,

    -- Limits
    max_usage_per_day INT,
    max_usage_per_month INT,
    max_usage_per_year INT,
    max_total_usage INT,
    max_discount_per_year BIGINT,
    usage_count INT DEFAULT 0,
    total_discount_given BIGINT DEFAULT 0,

    -- Reason and approval
    reason VARCHAR(200) NOT NULL,
    granted_by BIGINT,
    granted_at TIMESTAMP,
    approved_by BIGINT,
    approved_at TIMESTAMP,
    approval_notes TEXT,

    -- Validity period
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,

    -- Special conditions
    requires_approval BOOLEAN DEFAULT FALSE,
    requires_document BOOLEAN DEFAULT FALSE,
    document_path VARCHAR(500),
    conditions TEXT,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_suspended BOOLEAN DEFAULT FALSE,
    suspended_at TIMESTAMP,
    suspend_reason TEXT,

    -- Notes
    notes TEXT,
    internal_notes TEXT,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_emp_discount_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_emp_discount_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_emp_discount_item FOREIGN KEY (apply_to_item_id) REFERENCES items(id) ON DELETE SET NULL,
    CONSTRAINT fk_emp_discount_category FOREIGN KEY (apply_to_category_id) REFERENCES item_categories(id) ON DELETE SET NULL,
    CONSTRAINT fk_emp_discount_group FOREIGN KEY (apply_to_group_id) REFERENCES item_groups(id) ON DELETE SET NULL,
    CONSTRAINT fk_emp_discount_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_emp_discount_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_emp_discount_tenant ON employee_special_discounts(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_employee ON employee_special_discounts(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_code ON employee_special_discounts(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_type ON employee_special_discounts(discount_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_item ON employee_special_discounts(apply_to_item_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_category ON employee_special_discounts(apply_to_category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_group ON employee_special_discounts(apply_to_group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_priority ON employee_special_discounts(priority) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_start_date ON employee_special_discounts(start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_end_date ON employee_special_discounts(end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_is_active ON employee_special_discounts(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_discount_is_suspended ON employee_special_discounts(is_suspended) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE employee_special_discounts IS 'Individual discount rules for specific employees';
