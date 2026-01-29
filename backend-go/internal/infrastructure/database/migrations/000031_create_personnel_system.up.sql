-- Personnel System (Based on Refah/Yii structure)

-- 1. Relation Types (نسبت‌های خانوادگی)
CREATE TABLE IF NOT EXISTS relation_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    title_en VARCHAR(255),
    description TEXT,
    code_number INT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Insert default relation types (from stored procedure)
INSERT INTO relation_types (id, code, title, title_en, code_number) VALUES
(1, 'SPOUSE_FEMALE', 'همسر (زن)', 'Spouse (Female)', 1),
(2, 'SPOUSE_MALE', 'همسر (مرد)', 'Spouse (Male)', 2),
(3, 'CHILD', 'فرزند', 'Child', 3),
(4, 'DAUGHTER', 'دختر', 'Daughter', 4),
(5, 'SON', 'پسر', 'Son', 5),
(6, 'MOTHER', 'مادر', 'Mother', 6),
(7, 'FATHER', 'پدر', 'Father', 7),
(8, 'SELF', 'خود فرد (کارمند اصلی)', 'Self (Main Employee)', 8),
(11, 'OTHER', 'سایر', 'Other', 11)
ON CONFLICT (id) DO NOTHING;

-- 2. Guardianship Types (انواع کفالت)
CREATE TABLE IF NOT EXISTS guardianship_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    title_en VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Special Employee Types (گروه‌های ایثارگری)
CREATE TABLE IF NOT EXISTS special_employee_types (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    title_en VARCHAR(255),
    description TEXT,
    priority INT DEFAULT 0,
    discount_percentage DECIMAL(5,2),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Insert isar groups (from stored procedure logic)
INSERT INTO special_employee_types (id, code, title, title_en, priority) VALUES
(1, 'JANBAAZ_COMBINED', 'جانباز / رزمنده / ترکیبی', 'Veteran / Combined', 1),
(2, 'AZADEH', 'آزاده', 'Released POW', 2),
(3, 'SHAHID_CHILD_50', 'فرزند شاهد (50% جانبازی)', 'Martyr Child (50% Disability)', 3)
ON CONFLICT (id) DO NOTHING;

-- 4. Employees Table (کارمندان و افراد تحت تکفل)
CREATE TABLE IF NOT EXISTS employees (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Parent relationship (for family members)
    parent_id BIGINT REFERENCES employees(id) ON DELETE SET NULL,
    relation_type_id INT REFERENCES relation_types(id),

    -- Codes & Types
    custom_employee_code_id INT REFERENCES custom_employee_codes(id),
    special_employee_type_id INT REFERENCES special_employee_types(id),
    guardianship_type_id INT REFERENCES guardianship_types(id),

    -- Personal Info
    personnel_code VARCHAR(50),
    national_code VARCHAR(10),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    father_name VARCHAR(255),
    birth_date DATE,
    gender VARCHAR(10),
    marital_status VARCHAR(20),
    id_number VARCHAR(50),

    -- Contact
    phone VARCHAR(20),
    mobile VARCHAR(20),
    email VARCHAR(255),
    address TEXT,

    -- Employment Info
    branch_id INT,
    location_id INT,
    work_location_id INT,
    account_number VARCHAR(50),
    recruitment_date DATE,
    termination_date DATE,
    retirement_date DATE,

    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 1,
    status VARCHAR(20) DEFAULT 'active',

    -- Extra
    picture VARCHAR(255),
    description TEXT,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT fk_employee_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

-- 5. Employees Import Temp Table (برای sync از سرور HR)
CREATE TABLE IF NOT EXISTS employees_import_temp (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    -- Same structure as employees
    parent_id BIGINT,
    relation_type_id INT,
    custom_employee_code_id INT,
    special_employee_type_id INT,
    guardianship_type_id INT,

    personnel_code VARCHAR(50),
    national_code VARCHAR(10),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    father_name VARCHAR(255),
    birth_date DATE,
    gender VARCHAR(10),
    marital_status VARCHAR(20),
    id_number VARCHAR(50),

    phone VARCHAR(20),
    mobile VARCHAR(20),
    email VARCHAR(255),
    address TEXT,

    branch_id INT,
    location_id INT,
    work_location_id INT,
    account_number VARCHAR(50),
    recruitment_date DATE,
    termination_date DATE,
    retirement_date DATE,

    is_active BOOLEAN DEFAULT TRUE,
    priority INT DEFAULT 1,
    status VARCHAR(20),
    picture VARCHAR(255),
    description TEXT,

    -- Import metadata
    import_batch_id VARCHAR(100),
    import_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 6. Import History (تاریخچه sync)
CREATE TABLE IF NOT EXISTS employee_import_history (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    batch_id VARCHAR(100) UNIQUE NOT NULL,
    import_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    source VARCHAR(100),
    total_records INT DEFAULT 0,
    new_records INT DEFAULT 0,
    updated_records INT DEFAULT 0,
    failed_records INT DEFAULT 0,
    status VARCHAR(50) DEFAULT 'pending',
    notes TEXT,
    imported_by_user_id BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_employees_tenant ON employees(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_parent ON employees(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_personnel_code ON employees(personnel_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_national_code ON employees(national_code) WHERE deleted_at IS NULL;
CREATE INDEX idx_employees_relation_type ON employees(relation_type_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_employees_temp_batch ON employees_import_temp(import_batch_id);
CREATE INDEX idx_employees_temp_tenant ON employees_import_temp(tenant_id);

-- Comments
COMMENT ON TABLE employees IS 'کارمندان و افراد تحت تکفل - Compatible with Refah/Yii structure';
COMMENT ON TABLE employees_import_temp IS 'جدول موقت برای import از سرور منابع انسانی';
COMMENT ON COLUMN employees.parent_id IS 'برای افراد تحت تکفل - اشاره به کارمند اصلی';
COMMENT ON COLUMN employees.relation_type_id IS 'نسبت خانوادگی (همسر، فرزند و...)';
