-- Update claim_items table with new fields
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS item_id BIGINT;
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS prescription_item_id BIGINT;
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS instruction_id BIGINT;
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS dosage VARCHAR(100);
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS frequency VARCHAR(100);
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS duration VARCHAR(100);
ALTER TABLE claim_items ADD COLUMN IF NOT EXISTS body_site_id BIGINT;

-- Add foreign key constraints
ALTER TABLE claim_items ADD CONSTRAINT fk_claim_item_item FOREIGN KEY (item_id) REFERENCES items(id) ON DELETE SET NULL;
ALTER TABLE claim_items ADD CONSTRAINT fk_claim_item_prescription_item FOREIGN KEY (prescription_item_id) REFERENCES prescription_items(id) ON DELETE SET NULL;
ALTER TABLE claim_items ADD CONSTRAINT fk_claim_item_instruction FOREIGN KEY (instruction_id) REFERENCES instructions(id) ON DELETE SET NULL;
ALTER TABLE claim_items ADD CONSTRAINT fk_claim_item_body_site FOREIGN KEY (body_site_id) REFERENCES body_sites(id) ON DELETE SET NULL;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_claim_item_item ON claim_items(item_id);
CREATE INDEX IF NOT EXISTS idx_claim_item_prescription_item ON claim_items(prescription_item_id);
CREATE INDEX IF NOT EXISTS idx_claim_item_instruction ON claim_items(instruction_id);
CREATE INDEX IF NOT EXISTS idx_claim_item_body_site ON claim_items(body_site_id);

-- Update body_sites table with new fields
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS parent_id BIGINT;
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS category VARCHAR(50);
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS side VARCHAR(10);
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS icd10_code VARCHAR(20);
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS snomed_code VARCHAR(20);
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS cpt_modifier VARCHAR(10);
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE;
ALTER TABLE body_sites ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;

-- Make code unique if not already
CREATE UNIQUE INDEX IF NOT EXISTS uq_body_sites_code ON body_sites(code) WHERE deleted_at IS NULL;

-- Add indexes for body_sites
CREATE INDEX IF NOT EXISTS idx_body_sites_parent ON body_sites(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_body_sites_category ON body_sites(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_body_sites_icd10 ON body_sites(icd10_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_body_sites_snomed ON body_sites(snomed_code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_body_sites_is_active ON body_sites(is_active) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_body_sites_title_fa ON body_sites(title_fa) WHERE deleted_at IS NULL;

-- Add foreign key for body_sites parent
ALTER TABLE body_sites ADD CONSTRAINT fk_body_site_parent FOREIGN KEY (parent_id) REFERENCES body_sites(id) ON DELETE SET NULL;

-- Update existing body_sites with default values
UPDATE body_sites SET
    is_active = TRUE,
    sort_order = 0,
    side = 'NONE'
WHERE is_active IS NULL;

-- Add comments
COMMENT ON COLUMN claim_items.item_id IS 'Universal Item reference (drugs/services unified)';
COMMENT ON COLUMN claim_items.prescription_item_id IS 'Reference to prescription item if claim was created from prescription';
COMMENT ON COLUMN claim_items.instruction_id IS 'Drug/service usage instruction';
COMMENT ON COLUMN claim_items.body_site_id IS 'Direct reference to body site (anatomical location)';
COMMENT ON COLUMN body_sites.parent_id IS 'Parent body site for hierarchical structure';
COMMENT ON COLUMN body_sites.category IS 'Body site category (HEAD, TRUNK, UPPER_LIMB, LOWER_LIMB, etc.)';
COMMENT ON COLUMN body_sites.side IS 'Laterality (LEFT, RIGHT, BILATERAL, NONE)';
