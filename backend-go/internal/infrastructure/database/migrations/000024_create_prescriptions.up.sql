-- Create prescriptions table
CREATE TABLE IF NOT EXISTS prescriptions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Prescription number
    prescription_number VARCHAR(50) NOT NULL,

    -- Patient information
    employee_id BIGINT NOT NULL,
    is_for_main BOOLEAN DEFAULT TRUE,
    dependent_id BIGINT,

    -- Physician information
    physician_name VARCHAR(200) NOT NULL,
    physician_national_code VARCHAR(10),
    physician_medical_code VARCHAR(20),
    physician_specialty VARCHAR(100),
    physician_phone VARCHAR(20),

    -- Prescription details
    prescription_date TIMESTAMP NOT NULL,
    prescription_type VARCHAR(50) NOT NULL,
    is_electronic BOOLEAN DEFAULT FALSE,
    is_urgent BOOLEAN DEFAULT FALSE,

    -- Diagnosis
    main_diagnosis_code VARCHAR(20),
    main_diagnosis_title VARCHAR(300),
    secondary_diagnosis_codes TEXT,
    diagnosis_notes TEXT,

    -- Center/Facility
    center_id BIGINT,
    center_name VARCHAR(200),
    center_type VARCHAR(50),

    -- Approval and validation
    requires_approval BOOLEAN DEFAULT FALSE,
    approved_by_physician BOOLEAN DEFAULT FALSE,
    approval_date TIMESTAMP,
    approved_by BIGINT,
    approval_notes TEXT,

    -- Insurance validation
    is_validated_by_insurer BOOLEAN DEFAULT FALSE,
    validation_date TIMESTAMP,
    validated_by BIGINT,
    validation_notes TEXT,
    rejection_reason TEXT,

    -- Electronic prescription integration
    tamin_prescription_id VARCHAR(100),
    sepas_reference_id VARCHAR(100),
    national_rx_number VARCHAR(100),

    -- Files and attachments
    scanned_image_path VARCHAR(500),
    physician_signature_path VARCHAR(500),
    additional_documents_path TEXT,

    -- Conversion to claim
    is_converted_to_claim BOOLEAN DEFAULT FALSE,
    claim_id BIGINT,
    conversion_date TIMESTAMP,

    -- Status tracking
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    status_reason TEXT,
    submitted_at TIMESTAMP,
    completed_at TIMESTAMP,

    -- Notes
    physician_notes TEXT,
    pharmacy_notes TEXT,
    insurer_notes TEXT,
    internal_notes TEXT,

    -- Metadata
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_prescription_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_prescription_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
    CONSTRAINT fk_prescription_dependent FOREIGN KEY (dependent_id) REFERENCES family_members(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_center FOREIGN KEY (center_id) REFERENCES centers(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_claim FOREIGN KEY (claim_id) REFERENCES claims(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_approved_by FOREIGN KEY (approved_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_validated_by FOREIGN KEY (validated_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT uq_prescription_number_tenant UNIQUE (prescription_number, tenant_id, deleted_at)
);

-- Create prescription_items table
CREATE TABLE IF NOT EXISTS prescription_items (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    prescription_id BIGINT NOT NULL,

    -- Item reference
    item_id BIGINT,
    item_type VARCHAR(20) NOT NULL,

    -- Manual entry
    manual_item_name VARCHAR(300),
    manual_item_code VARCHAR(50),

    -- Prescription details
    quantity INT NOT NULL,
    dosage VARCHAR(100),
    frequency VARCHAR(100),
    duration VARCHAR(100),
    route VARCHAR(50),
    instruction_id BIGINT,
    special_notes TEXT,

    -- Body site
    body_site_id BIGINT,

    -- Priority and urgency
    is_urgent BOOLEAN DEFAULT FALSE,
    priority INT DEFAULT 0,
    sort_order INT DEFAULT 0,

    -- Validation
    is_validated BOOLEAN DEFAULT FALSE,
    validation_notes TEXT,
    substitution_allowed BOOLEAN DEFAULT TRUE,

    -- Status
    status VARCHAR(20) DEFAULT 'PENDING',
    is_active BOOLEAN DEFAULT TRUE,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_prescription_item_prescription FOREIGN KEY (prescription_id) REFERENCES prescriptions(id) ON DELETE CASCADE,
    CONSTRAINT fk_prescription_item_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_item_instruction FOREIGN KEY (instruction_id) REFERENCES instructions(id) ON DELETE SET NULL,
    CONSTRAINT fk_prescription_item_body_site FOREIGN KEY (body_site_id) REFERENCES body_sites(id) ON DELETE SET NULL
);

-- Create indexes for prescriptions
CREATE INDEX idx_prescription_tenant ON prescriptions(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_employee ON prescriptions(employee_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_dependent ON prescriptions(dependent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_date ON prescriptions(prescription_date) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_type ON prescriptions(prescription_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_is_electronic ON prescriptions(is_electronic) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_is_for_main ON prescriptions(is_for_main) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_physician_code ON prescriptions(physician_medical_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_main_diagnosis ON prescriptions(main_diagnosis_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_is_validated ON prescriptions(is_validated_by_insurer) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_tamin_id ON prescriptions(tamin_prescription_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_sepas_id ON prescriptions(sepas_reference_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_national_rx ON prescriptions(national_rx_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_is_converted ON prescriptions(is_converted_to_claim) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_claim_id ON prescriptions(claim_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_status ON prescriptions(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_is_active ON prescriptions(is_active) WHERE deleted_at IS NULL;

-- Create indexes for prescription_items
CREATE INDEX idx_prescription_item_tenant ON prescription_items(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_item_prescription ON prescription_items(prescription_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_item_item ON prescription_items(item_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_prescription_item_body_site ON prescription_items(body_site_id) WHERE deleted_at IS NULL;

-- Add comments
COMMENT ON TABLE prescriptions IS 'Medical prescriptions - can be converted to claims';
COMMENT ON TABLE prescription_items IS 'Items in prescriptions (drugs/services)';
