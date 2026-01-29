-- Create employee_illnesses table
CREATE TABLE IF NOT EXISTS employee_illnesses (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,

    -- Illness information
    icd10_code VARCHAR(20) NOT NULL,
    icd10_title VARCHAR(300) NOT NULL,
    icd10_title_fa VARCHAR(300) NOT NULL,
    icd10_category VARCHAR(50),

    -- Diagnosis details
    diagnosis_date TIMESTAMP,
    diagnosed_by VARCHAR(200),
    diagnosis_center VARCHAR(200),
    diagnosis_details TEXT,

    -- Severity and status
    severity VARCHAR(20),
    status VARCHAR(20),
    is_chronic BOOLEAN DEFAULT FALSE,
    is_congenital BOOLEAN DEFAULT FALSE,
    is_hereditary BOOLEAN DEFAULT FALSE,

    -- Insurance impact
    is_pre_existing BOOLEAN DEFAULT FALSE,
    requires_declaration BOOLEAN DEFAULT TRUE,
    affects_premium BOOLEAN DEFAULT FALSE,
    premium_loading_percent DECIMAL(5,2),

    -- Coverage rules
    is_covered BOOLEAN DEFAULT TRUE,
    coverage_start_date TIMESTAMP,
    coverage_end_date TIMESTAMP,
    waiting_period_days INT,
    exclusion_reason TEXT,
    special_conditions TEXT,
    annual_coverage_limit BIGINT,

    -- Treatment information
    current_treatment TEXT,
    requires_medication BOOLEAN DEFAULT FALSE,
    requires_monitoring BOOLEAN DEFAULT FALSE,
    last_visit_date TIMESTAMP,
    next_visit_date TIMESTAMP,

    -- Related documents
    medical_report_path VARCHAR(500),
    lab_test_results_path VARCHAR(500),
    imaging_results_path VARCHAR(500),

    -- Notes and remarks
    physician_notes TEXT,
    insurance_notes TEXT,
    internal_notes TEXT,

    -- Approval workflow
    declared_by_employee BOOLEAN DEFAULT FALSE,
    declaration_date TIMESTAMP,
    verified_by_insurer BOOLEAN DEFAULT FALSE,
    verification_date TIMESTAMP,
    verified_by BIGINT,
    rejection_reason TEXT,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_emp_illness_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_emp_illness_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_emp_illness_verified_by FOREIGN KEY (verified_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_emp_illness_tenant ON employee_illnesses(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_employee ON employee_illnesses(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_icd10 ON employee_illnesses(icd10_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_category ON employee_illnesses(icd10_category) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_diagnosis_date ON employee_illnesses(diagnosis_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_severity ON employee_illnesses(severity) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_status ON employee_illnesses(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_chronic ON employee_illnesses(is_chronic) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_congenital ON employee_illnesses(is_congenital) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_hereditary ON employee_illnesses(is_hereditary) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_pre_existing ON employee_illnesses(is_pre_existing) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_covered ON employee_illnesses(is_covered) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_coverage_start ON employee_illnesses(coverage_start_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_emp_illness_is_active ON employee_illnesses(is_active) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE employee_illnesses IS 'Chronic illnesses and pre-existing conditions of employees';
