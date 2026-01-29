-- Create instructions table
CREATE TABLE IF NOT EXISTS instructions (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Code and naming
    code VARCHAR(50) NOT NULL,
    title VARCHAR(200) NOT NULL,
    title_fa VARCHAR(200) NOT NULL,

    -- Instruction details
    frequency VARCHAR(100),
    duration VARCHAR(100),
    timing VARCHAR(100),
    route VARCHAR(100),
    dosage VARCHAR(100),
    description TEXT,

    -- Templates
    is_template BOOLEAN DEFAULT FALSE,
    template TEXT,

    -- Category
    category VARCHAR(50),

    -- Special flags
    requires_supervision BOOLEAN DEFAULT FALSE,
    is_emergency BOOLEAN DEFAULT FALSE,
    is_chronic BOOLEAN DEFAULT FALSE,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,

    -- Metadata
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_instruction_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT uq_instruction_code_tenant UNIQUE (code, tenant_id, deleted_at)
);

-- Create indexes
CREATE INDEX idx_instruction_tenant ON instructions(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_instruction_title_fa ON instructions(title_fa) WHERE deleted_at IS NULL;
CREATE INDEX idx_instruction_category ON instructions(category) WHERE deleted_at IS NULL;
CREATE INDEX idx_instruction_is_template ON instructions(is_template) WHERE deleted_at IS NULL;
CREATE INDEX idx_instruction_is_active ON instructions(is_active) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON TABLE instructions IS 'Drug/service usage instructions for prescriptions';

-- Insert seed data (common instructions)
INSERT INTO instructions (tenant_id, code, title, title_fa, frequency, timing, route, category, is_template, is_active, sort_order) VALUES
(1, 'BID_ORAL', 'Twice daily oral', 'روزی دو بار خوراکی', 'روزی 2 بار', 'با غذا', 'خوراکی', 'DRUG', true, true, 1),
(1, 'TID_ORAL', 'Three times daily oral', 'روزی سه بار خوراکی', 'روزی 3 بار', 'بعد از غذا', 'خوراکی', 'DRUG', true, true, 2),
(1, 'QID_ORAL', 'Four times daily oral', 'روزی چهار بار خوراکی', 'روزی 4 بار', 'بعد از غذا', 'خوراکی', 'DRUG', true, true, 3),
(1, 'QD_ORAL', 'Once daily oral', 'روزی یک بار خوراکی', 'روزی 1 بار', 'قبل از خواب', 'خوراکی', 'DRUG', true, true, 4),
(1, 'PRN_ORAL', 'As needed oral', 'در صورت نیاز خوراکی', 'در صورت نیاز', 'هر زمان', 'خوراکی', 'DRUG', true, true, 5),
(1, 'IM_INJECTION', 'Intramuscular', 'تزریق عضلانی', 'طبق دستور پزشک', '', 'تزریقی - عضلانی', 'DRUG', true, true, 6),
(1, 'IV_INJECTION', 'Intravenous', 'تزریق وریدی', 'طبق دستور پزشک', '', 'تزریقی - وریدی', 'DRUG', true, true, 7),
(1, 'TOPICAL', 'Topical application', 'مصرف موضعی', 'روزی 2-3 بار', '', 'موضعی', 'DRUG', true, true, 8),
(1, 'EYE_DROP', 'Eye drops', 'قطره چشمی', 'روزی 3 بار', '', 'چشمی', 'DRUG', true, true, 9),
(1, 'EAR_DROP', 'Ear drops', 'قطره گوشی', 'روزی 2 بار', '', 'گوشی', 'DRUG', true, true, 10);
