# Architecture Guardrails - TPA System

## Overview

This document defines the architectural constraints and guardrails that must be followed when developing the TPA system. These rules ensure consistency, maintainability, and proper separation of concerns.

## 1. Module Boundaries

### Go Backend (TPA Core)
Responsible for:
- Claims processing (`/api/v1/claims/*`)
- Package management (`/api/v1/packages/*`)
- Provider/Center management (`/api/v1/centers/*`)
- Financial settlements (`/api/v1/settlements/*`)
- Drug & Service catalog (`/api/v1/drugs/*`, `/api/v1/services/*`)
- Rule engine evaluation
- External integrations (Tamin, Sepas, IRC)

**MUST NOT** handle:
- User authentication (delegates to NestJS or shared JWT)
- Commission case management
- Social work cases

### NestJS Backend (Commission Module)
Responsible for:
- User management (`/api/v1/users/*`)
- Authentication (`/api/v1/auth/*`)
- Medical commission cases (`/api/v1/cases/*`)
- Social work cases (`/api/v1/social-work/*`)
- Verdict management
- Insured person registry

**MUST NOT** handle:
- Claim processing
- Financial settlements
- Direct database access to Go-owned tables

## 2. Dependency Direction

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (Vue 3)                        │
│                    Depends on both APIs                      │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
     ┌────────────────┐              ┌────────────────┐
     │   Go Backend   │◄────────────►│ NestJS Backend │
     │   (TPA Core)   │   Events     │  (Commission)  │
     └────────────────┘              └────────────────┘
              │                               │
              └───────────────┬───────────────┘
                              ▼
                    ┌─────────────────┐
                    │   PostgreSQL    │
                    │  (Shared DB)    │
                    └─────────────────┘
```

### Rules:
1. **No direct cross-module database access** - Modules communicate via Events or APIs only
2. **Shared tables** are read-only for the non-owning module
3. **Events flow bidirectionally** between Go and NestJS
4. **Frontend** may call both backends through nginx proxy

## 3. Multi-Tenancy Requirements

### Every Request MUST:
1. Include tenant identification (header `X-Tenant-ID` or JWT claim)
2. Be validated through tenant middleware
3. Have tenant_id in query scope

### Database Operations MUST:
1. Include `tenant_id` in WHERE clauses
2. Never allow cross-tenant data access
3. Use tenant-scoped connections or middleware

### Tenant Isolation Levels:
```
Level 1: Row-level (tenant_id column) - DEFAULT
Level 2: Schema-level (tenant_N schema) - For large tenants
Level 3: Database-level - For enterprise clients
```

### Code Pattern (Go):
```go
// CORRECT - Always scope by tenant
func (r *ClaimRepo) GetByID(ctx context.Context, id uint) (*Claim, error) {
    tenantID := tenant.FromContext(ctx)
    var claim Claim
    err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&claim).Error
    return &claim, err
}

// INCORRECT - Never access without tenant scope
func (r *ClaimRepo) GetByID(ctx context.Context, id uint) (*Claim, error) {
    var claim Claim
    err := r.db.Where("id = ?", id).First(&claim).Error  // VIOLATION!
    return &claim, err
}
```

## 4. Event-Driven Communication

### Event Categories:
1. **Domain Events** - Business state changes (CommissionVerdictIssued)
2. **Integration Events** - Cross-module communication
3. **System Events** - Infrastructure notifications

### Event Rules:
1. Events are **immutable** - Never modify published events
2. Events are **versioned** - Include schema version
3. Events are **idempotent** - Handlers must handle duplicates
4. Events include **correlation_id** for tracing

### Event Schema (Required Fields):
```json
{
  "event_id": "uuid",
  "event_type": "domain.entity.action",
  "version": "1.0.0",
  "timestamp": "ISO8601",
  "source": "service-name",
  "tenant_id": 123,
  "correlation_id": "uuid",
  "causation_id": "uuid"
}
```

### Event Flow Example:
```
Commission Module                    TPA Core
      │                                 │
      │  commission.verdict.issued      │
      │────────────────────────────────►│
      │                                 │
      │                                 │ Process financial impact
      │                                 │
      │   tpa.coverage.updated          │
      │◄────────────────────────────────│
      │                                 │
```

## 5. Database Schema Conventions

### Naming:
- Tables: `snake_case` plural (`claims`, `claim_items`)
- Columns: `snake_case` (`tenant_id`, `created_at`)
- Indexes: `idx_{table}_{columns}` (`idx_claims_tenant_status`)
- Foreign keys: `fk_{table}_{referenced}` (`fk_claims_center`)

### Required Columns (All Tables):
```sql
id          BIGSERIAL PRIMARY KEY
tenant_id   BIGINT NOT NULL REFERENCES tenants(id)
created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
deleted_at  TIMESTAMP WITH TIME ZONE  -- Soft delete
```

### Table Ownership:
| Module  | Tables |
|---------|--------|
| Go      | claims, claim_items, packages, settlements, centers, drugs, services, tariffs, rules, rule_execution_logs |
| NestJS  | commission_users, commission_cases, commission_verdicts, social_work_cases, commission_provinces |
| Shared  | tenants, persons, employees, policies (Go writes, both read) |

### Cross-Module Access:
```sql
-- NestJS reading Go-owned table (READ ONLY)
SELECT id, full_name FROM persons WHERE national_id = '1234567890';

-- NestJS writing to own table
INSERT INTO commission_cases (...) VALUES (...);

-- NEVER: NestJS writing to Go-owned table
-- INSERT INTO claims (...) VALUES (...);  -- VIOLATION!
```

## 6. API Versioning

### URL Structure:
```
/api/v{major}/resource
```

### Version Lifecycle:
1. **Current** (v1) - Active development
2. **Supported** (v0) - Bug fixes only
3. **Deprecated** - 6 month notice before removal
4. **Removed** - Returns 410 Gone

### Breaking Changes Require:
1. New major version
2. Deprecation notice on old version
3. Migration guide documentation

### Non-Breaking Changes (Same Version):
- Adding optional fields
- Adding new endpoints
- Adding new enum values (if clients handle unknown)

## 7. Security Guardrails

### Authentication:
1. All endpoints require authentication except `/auth/login`, `/health`
2. JWT tokens expire in 24 hours
3. Refresh tokens expire in 7 days

### Authorization:
1. Role-based access control (RBAC)
2. Tenant isolation enforced at middleware level
3. Resource-level permissions for sensitive operations

### Data Protection:
1. PII fields encrypted at rest
2. Passwords hashed with bcrypt (cost 12+)
3. Sensitive data masked in logs

### Input Validation:
1. All inputs validated before processing
2. SQL injection prevention via parameterized queries
3. XSS prevention via output encoding

## 8. Error Handling

### Error Response Format:
```json
{
  "error": {
    "code": "CLAIM_NOT_FOUND",
    "message": "Claim with ID 123 not found",
    "details": {},
    "trace_id": "uuid"
  }
}
```

### HTTP Status Codes:
| Code | Usage |
|------|-------|
| 200  | Success |
| 201  | Created |
| 400  | Validation error |
| 401  | Not authenticated |
| 403  | Not authorized |
| 404  | Resource not found |
| 409  | Conflict (duplicate) |
| 422  | Business rule violation |
| 500  | Internal error |

### Error Logging:
1. 4xx errors - INFO level
2. 5xx errors - ERROR level with stack trace
3. Include correlation_id in all logs

## 9. Testing Requirements

### Unit Tests:
- Domain logic: 80%+ coverage
- Use cases: 70%+ coverage
- Repositories: Integration tests preferred

### Integration Tests:
- All API endpoints
- Database migrations
- Event handlers

### E2E Tests:
- Critical user flows
- Cross-module communication

## 10. Deployment Guardrails

### Pre-Deployment Checklist:
- [ ] All tests pass
- [ ] Database migrations tested
- [ ] Environment variables documented
- [ ] Rollback plan prepared

### Zero-Downtime Deployment:
1. Database migrations must be backward compatible
2. New code must handle old data format
3. Feature flags for risky changes

### Health Checks:
```
GET /health          - Basic health
GET /health/ready    - Ready to serve traffic
GET /health/live     - Process is alive
```

## Enforcement

These guardrails are enforced through:
1. **Code Review** - Required for all PRs
2. **Linting** - Automated checks in CI
3. **Architecture Tests** - Dependency validation
4. **Security Scans** - Vulnerability detection

## Changelog

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2024-01 | Initial guardrails document |
