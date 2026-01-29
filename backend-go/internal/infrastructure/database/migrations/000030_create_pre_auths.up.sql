-- Create pre_auths table (Pre-Authorization / علی‌الحساب)
CREATE TABLE IF NOT EXISTS pre_auths (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,

    person_id BIGINT NOT NULL,
    subject VARCHAR(500),
    amount BIGINT NOT NULL,
    type SMALLINT,

    payment_date TIMESTAMP,

    register_user_id BIGINT NOT NULL,
    register_date TIMESTAMP NOT NULL,

    claim_id BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_pre_auth_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_pre_auth_person FOREIGN KEY (person_id) REFERENCES persons(id) ON DELETE CASCADE
);

CREATE INDEX idx_pre_auths_tenant ON pre_auths(tenant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pre_auths_person ON pre_auths(person_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_pre_auths_claim ON pre_auths(claim_id) WHERE deleted_at IS NULL;

COMMENT ON TABLE pre_auths IS 'Pre-authorization (علی‌الحساب) - advance payments';
