-- Remove foreign keys from claim_items
ALTER TABLE claim_items DROP CONSTRAINT IF EXISTS fk_claim_item_item;
ALTER TABLE claim_items DROP CONSTRAINT IF EXISTS fk_claim_item_prescription_item;
ALTER TABLE claim_items DROP CONSTRAINT IF EXISTS fk_claim_item_instruction;
ALTER TABLE claim_items DROP CONSTRAINT IF EXISTS fk_claim_item_body_site;

-- Remove indexes from claim_items
DROP INDEX IF EXISTS idx_claim_item_item;
DROP INDEX IF EXISTS idx_claim_item_prescription_item;
DROP INDEX IF EXISTS idx_claim_item_instruction;
DROP INDEX IF EXISTS idx_claim_item_body_site;

-- Remove columns from claim_items
ALTER TABLE claim_items DROP COLUMN IF EXISTS item_id;
ALTER TABLE claim_items DROP COLUMN IF EXISTS prescription_item_id;
ALTER TABLE claim_items DROP COLUMN IF EXISTS instruction_id;
ALTER TABLE claim_items DROP COLUMN IF EXISTS dosage;
ALTER TABLE claim_items DROP COLUMN IF EXISTS frequency;
ALTER TABLE claim_items DROP COLUMN IF EXISTS duration;
ALTER TABLE claim_items DROP COLUMN IF EXISTS body_site_id;

-- Remove foreign key from body_sites
ALTER TABLE body_sites DROP CONSTRAINT IF EXISTS fk_body_site_parent;

-- Remove indexes from body_sites
DROP INDEX IF EXISTS uq_body_sites_code;
DROP INDEX IF EXISTS idx_body_sites_parent;
DROP INDEX IF EXISTS idx_body_sites_category;
DROP INDEX IF EXISTS idx_body_sites_icd10;
DROP INDEX IF EXISTS idx_body_sites_snomed;
DROP INDEX IF EXISTS idx_body_sites_is_active;
DROP INDEX IF EXISTS idx_body_sites_title_fa;

-- Remove columns from body_sites
ALTER TABLE body_sites DROP COLUMN IF EXISTS parent_id;
ALTER TABLE body_sites DROP COLUMN IF EXISTS category;
ALTER TABLE body_sites DROP COLUMN IF EXISTS side;
ALTER TABLE body_sites DROP COLUMN IF EXISTS icd10_code;
ALTER TABLE body_sites DROP COLUMN IF EXISTS snomed_code;
ALTER TABLE body_sites DROP COLUMN IF EXISTS cpt_modifier;
ALTER TABLE body_sites DROP COLUMN IF EXISTS description;
ALTER TABLE body_sites DROP COLUMN IF EXISTS is_active;
ALTER TABLE body_sites DROP COLUMN IF EXISTS sort_order;
