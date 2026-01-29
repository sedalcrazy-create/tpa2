-- Create insurance_histories table
CREATE TABLE IF NOT EXISTS insurance_histories (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,
    insurance_id BIGINT NOT NULL,
    policy_id BIGINT,

    -- History event type
    event_type VARCHAR(50) NOT NULL,

    -- Period
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,

    -- Coverage details snapshot
    coverage_type VARCHAR(50),
    plan_name VARCHAR(200),
    plan_code VARCHAR(50),
    annual_limit BIGINT,
    franchise_amount BIGINT,
    franchise_percentage DECIMAL(5,2),
    premium_amount BIGINT,

    -- Change details
    change_reason TEXT,
    changed_fields TEXT,
    previous_values TEXT,
    new_values TEXT,

    -- Suspension details
    suspension_reason TEXT,
    suspended_at TIMESTAMP,
    suspended_by BIGINT,

    -- Termination details
    termination_reason TEXT,
    terminated_at TIMESTAMP,
    terminated_by BIGINT,
    refund_amount BIGINT,

    -- Renewal details
    renewal_date TIMESTAMP,
    renewal_type VARCHAR(50),
    premium_change_percent DECIMAL(5,2),

    -- Usage statistics
    total_claims INT DEFAULT 0,
    total_claim_amount BIGINT DEFAULT 0,
    total_paid_amount BIGINT DEFAULT 0,
    loss_ratio DECIMAL(5,2),

    -- Documents
    contract_path VARCHAR(500),
    addendum_path VARCHAR(500),
    termination_letter_path VARCHAR(500),
    documents_path TEXT,

    -- Approval workflow
    requires_approval BOOLEAN DEFAULT FALSE,
    approved_by BIGINT,
    approved_at TIMESTAMP,
    approval_notes TEXT,
    rejected_by BIGINT,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,

    -- Notes
    description TEXT,
    notes TEXT,
    internal_notes TEXT,

    -- System tracking
    created_by BIGINT,
    is_system_generated BOOLEAN DEFAULT FALSE,
    source_system VARCHAR(50),

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    is_current BOOLEAN DEFAULT FALSE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_ins_history_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_ins_history_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_ins_history_insurance FOREIGN KEY (insurance_id) REFERENCES insurances(id) ON DELETE CASCADE,
    CONSTRAINT fk_ins_history_policy FOREIGN KEY (policy_id) REFERENCES policies(id) ON DELETE SET NULL,
    CONSTRAINT fk_ins_history_suspended_by FOREIGN KEY (suspended_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_ins_history_terminated_by FOREIGN KEY (terminated_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_ins_history_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_ins_history_rejected_by FOREIGN KEY (rejected_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_ins_history_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_ins_history_tenant ON insurance_histories(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_employee ON insurance_histories(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_insurance ON insurance_histories(insurance_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_policy ON insurance_histories(policy_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_event_type ON insurance_histories(event_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_start_date ON insurance_histories(start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_end_date ON insurance_histories(end_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_plan_code ON insurance_histories(plan_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_is_active ON insurance_histories(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_ins_history_is_current ON insurance_histories(is_current) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE insurance_histories IS 'History of insurance coverage for employees - tracks all changes, suspensions, and renewals';
