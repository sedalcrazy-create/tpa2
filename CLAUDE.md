# CLAUDE.md - TPA System (سامانه جامع TPA درمانی)

## Project Overview

A comprehensive Third Party Administrator (TPA) system for health insurance management. Built with Go backend, supporting multiple insurers, claim processing, provider settlements, and integration with national health systems (Tamin, Sepas, IRC).

## Technology Stack

- **Backend**: Go 1.22+ with Fiber framework
- **Database**: PostgreSQL 16 (schema-based multi-tenancy)
- **Cache/Queue**: Redis
- **Rule Engine**: Grule (github.com/hyperjumptech/grule-rule-engine)
- **Frontend**: Vue 3 + Quasar (or keep existing)
- **Auth**: JWT + OAuth2 (Bale integration ready)
- **API Docs**: OpenAPI/Swagger

## Project Structure

```
/e/project/TPA2/
├── CLAUDE.md                    # This file
├── docker-compose.yml           # Unified deployment (one DB, multiple backends)
├── backend-go/                  # Go Backend - TPA Core (Claims, Centers, Settlements)
│   ├── cmd/
│   │   ├── api/                # Main API server
│   │   │   └── main.go
│   │   └── worker/             # Background job worker
│   │       └── main.go
│   ├── internal/
│   │   ├── config/             # Configuration
│   │   ├── domain/             # Domain entities & interfaces
│   │   │   ├── entity/         # Core entities
│   │   │   ├── repository/     # Repository interfaces
│   │   │   └── service/        # Domain services
│   │   ├── usecase/            # Application use cases
│   │   ├── infrastructure/     # Database & external services
│   │   ├── delivery/http/      # HTTP handlers & middleware
│   │   └── pkg/                # Internal packages
│   │       ├── grule/          # Rule engine wrapper
│   │       ├── tamin/          # Tamin API client
│   │       ├── sepas/          # Sepas API client
│   │       └── irc/            # IRC drug database client
│   ├── Dockerfile
│   └── go.mod
├── backend-nestjs/              # NestJS Backend - Commission & Social Work
│   ├── src/
│   │   ├── auth/               # Authentication (shared with TPA)
│   │   ├── users/              # User management
│   │   ├── cases/              # Medical commission cases
│   │   ├── social-work/        # Social work cases
│   │   ├── insured-persons/    # Insured persons management
│   │   ├── case-types/         # Case type definitions
│   │   ├── verdict-templates/  # Verdict templates
│   │   ├── provinces/          # Province management
│   │   ├── dashboard/          # Dashboard stats
│   │   ├── import/             # Data import
│   │   └── entities/           # TypeORM entities
│   ├── Dockerfile
│   └── package.json
├── frontend/                    # Vue 3 Unified Frontend
│   ├── src/
│   │   ├── views/              # All views (TPA + Commission + Social Work)
│   │   ├── layouts/            # MainLayout with unified menu
│   │   ├── stores/             # Pinia stores
│   │   ├── services/           # API services
│   │   └── router/             # Vue Router
│   ├── Dockerfile
│   └── nginx.conf              # Proxy to both backends
├── grule-rule-engine-master/   # Rule engine reference
└── docs/                        # Documentation
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                       Frontend (Vue 3)                          │
│         /tpa - Unified UI for all modules                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                      ┌───────┴───────┐
                      │   Nginx       │
                      │   (Proxy)     │
                      └───────┬───────┘
              ┌───────────────┴───────────────┐
              │                               │
     ┌────────┴────────┐            ┌────────┴────────┐
     │  Go Backend     │            │  NestJS Backend │
     │  (Port 8080)    │            │  (Port 3000)    │
     │                 │            │                 │
     │  - Claims       │            │  - Commission   │
     │  - Packages     │            │  - Social Work  │
     │  - Centers      │            │  - Users        │
     │  - Settlements  │            │  - Auth         │
     │  - Drugs        │            │  - Insured      │
     │  - Services     │            │  - Case Types   │
     └────────┬────────┘            └────────┬────────┘
              │                               │
              └───────────────┬───────────────┘
                              │
                    ┌─────────┴─────────┐
                    │   PostgreSQL      │
                    │   (Shared DB)     │
                    │   Database: tpa   │
                    └───────────────────┘
```

## Core Modules

### 1. Personnel (پرسنل) - از Refah
- Employee management (بیمه‌شده اصلی)
- Family/Dependents (افراد تبعی)
- Insurance coverage periods
- Special groups & discounts

### 2. Drug & Service Bank (بانک دارو و خدمات)
- Drug database (IRC integration)
- Service catalog with K factors
- Drug interactions
- Price management

### 3. Provider Management (مراکز درمانی)
- Center registration & contracts
- Payment accounts (Sheba)
- Accreditation levels
- Contract terms

### 4. Claim Processing (پردازش ادعا) - از Sinaps
- 12 claim types (drug, hospital, dental, lab, imaging, etc.)
- Claim items (drugs/services)
- Diagnosis codes
- Document attachments

### 5. Pricing Engine (موتور قیمت‌گذاری) - Grule
- Dynamic pricing rules
- Coverage limits
- Franchise calculation
- Deductions

### 6. Financial Settlement (تسویه مالی)
- Package management (بسته اسناد)
- Invoice aggregation
- Provider payments
- Reports

### 7. Integrations (یکپارچه‌سازی)
- Tamin API (نسخ الکترونیک)
- Sepas (پرونده سلامت)
- IRC (قیمت دارو)

## Database Schema (PostgreSQL)

```sql
-- Multi-tenant with schemas
CREATE SCHEMA shared;      -- Common data (provinces, drugs, services)
CREATE SCHEMA insurer_1;   -- Insurer-specific data
CREATE SCHEMA insurer_2;   -- Another insurer

-- Key tables per schema:
-- personnel, policies, claims, claim_items, centers, packages, settlements
```

## Key Entities

### Claim (ادعا)
```go
type Claim struct {
    ID                 uint
    PolicyMemberID     uint          // بیمه‌شده
    CenterID           uint          // مرکز درمانی
    PackageID          *uint         // بسته ارسالی
    Type               ClaimType     // نوع ادعا
    Status             ClaimStatus   // وضعیت
    AdmissionDate      time.Time     // تاریخ پذیرش
    DischargeDate      *time.Time    // تاریخ ترخیص
    RequestAmount      int64         // مبلغ درخواستی
    ApprovedAmount     int64         // مبلغ تایید شده
    Items              []ClaimItem   // اقلام
    Diagnoses          []Diagnosis   // تشخیص‌ها
}
```

### ClaimItem (قلم ادعا)
```go
type ClaimItem struct {
    ID                 uint
    ClaimID            uint
    DrugID             *uint         // دارو
    ServiceID          *uint         // خدمت
    Count              int           // تعداد
    RequestPrice       int64         // مبلغ درخواستی
    ApprovedPrice      int64         // تایید شده
    BasicInsShare      int64         // سهم بیمه پایه
    SupplementShare    int64         // سهم بیمه تکمیلی
    Franchise          int64         // فرانشیز
    Deduction          int64         // کسورات
}
```

### Center (مرکز درمانی)
```go
type Center struct {
    ID                 uint
    Title              string
    SiamID             string
    Type               CenterType
    ProvinceID         uint
    Level              int           // سطح اعتباربخشی
    ContractStatus     ContractStatus
    AccountNumber      string
    ShebaNumber        string
    PaymentID          string
}
```

### Package (بسته اسناد)
```go
type Package struct {
    ID                 uint
    CenterID           uint
    Claims             []Claim
    TotalAmount        int64
    Status             PackageStatus
    LetterNumber       string
    PaymentDate        *time.Time
    PaymentAmount      int64
}
```

## Claim Status Flow

```
WaitRegister → WaitCheck → WaitCheckConfirm → WaitSendToFinancial → Archive
                    ↓
                Returned (عودت)
```

## Claim Types (CompositionType)

1. Drug (داروخانه)
2. Hospitalization (بستری)
3. Dental (دندانپزشکی)
4. DoctorService (ویزیت)
5. ParaclinicalTest (آزمایشگاه)
6. ParaclinicalImage (تصویربرداری)
7. ParaclinicalPhysiotherapy (فیزیوتراپی)
8. OutpatientSurgery (جراحی سرپایی)
9. Emergency (اورژانس)
10. MedicalEquipment (تجهیزات)
11. Injections (تزریقات)
12. Clinic (سرپایی بیمارستان)

## API Endpoints (Planned)

```
# Auth
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh

# Personnel
GET    /api/v1/personnel
GET    /api/v1/personnel/:id
POST   /api/v1/personnel
PUT    /api/v1/personnel/:id
GET    /api/v1/personnel/:id/family

# Claims
GET    /api/v1/claims
GET    /api/v1/claims/:id
POST   /api/v1/claims
PUT    /api/v1/claims/:id
POST   /api/v1/claims/:id/items
POST   /api/v1/claims/:id/submit
POST   /api/v1/claims/:id/approve

# Centers
GET    /api/v1/centers
GET    /api/v1/centers/:id
POST   /api/v1/centers
GET    /api/v1/centers/:id/packages

# Packages
GET    /api/v1/packages
GET    /api/v1/packages/:id
POST   /api/v1/packages/:id/approve
POST   /api/v1/packages/:id/pay

# Integrations
POST   /api/v1/tamin/inquiry
POST   /api/v1/sepas/inquiry
GET    /api/v1/irc/drugs
```

## Running the Project

### Backend
```bash
cd backend-go
cp .env.example .env
# Edit .env with your settings

# Run migrations
make migrate-up

# Run API server
make run-api

# Run worker (background jobs)
make run-worker
```

### Database
```bash
# PostgreSQL in Docker
docker run -d \
  --name tpa-postgres \
  -e POSTGRES_USER=tpa \
  -e POSTGRES_PASSWORD=tpa123 \
  -e POSTGRES_DB=tpa \
  -p 5432:5432 \
  postgres:16

# Redis
docker run -d \
  --name tpa-redis \
  -p 6379:6379 \
  redis:7-alpine
```

## Development Guidelines

### Code Style
- Follow Go conventions (gofmt, golint)
- Use meaningful Persian comments for business logic
- Keep handlers thin, logic in usecases
- Use interfaces for testability

### Commit Messages
- feat: New feature
- fix: Bug fix
- refactor: Code refactoring
- docs: Documentation
- test: Tests

### Testing
```bash
make test          # Run all tests
make test-coverage # With coverage report
```

## Reference Projects

- **Sinaps (Jafam)**: `New folder (2)/Projects/` - .NET CQRS architecture, Claim/ClaimItem structure
- **Refah (Yii)**: `New folder (2)/refah/` - Personnel system, ItemPriceConditions, Invoice
- **Grule**: `grule-rule-engine-master/` - Rule engine for Go

## Environment Variables

```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=tpa
DB_PASSWORD=tpa123
DB_NAME=tpa

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRE=24h

# Tamin API
TAMIN_API_URL=https://soa.tamin.ir/
TAMIN_CLIENT_ID=xxx
TAMIN_SECRET=xxx

# Sepas API
SEPAS_API_URL=xxx
SEPAS_TOKEN=xxx
```

## TODO / Roadmap

### Phase 1: Core Setup ✅
- [x] Project structure
- [x] CLAUDE.md
- [x] Go module init (go.mod)
- [x] Database setup (GORM + PostgreSQL)
- [x] Base entities (base.go, enums.go)

### Phase 2: Personnel Module ✅
- [x] Person entity
- [x] Employee entity
- [x] FamilyMember entity
- [x] Policy/PolicyMember
- [ ] CRUD APIs (basic structure)

### Phase 3: Drug/Service Bank ✅
- [x] Drug entity (drug.go)
- [x] DrugPrice, DrugInteraction
- [x] Service entity (service.go)
- [x] ServicePrice, Tariff
- [ ] IRC integration

### Phase 4: Claim Processing ✅
- [x] Claim entity (claim.go)
- [x] ClaimItem entity
- [x] ClaimDiagnosis, ClaimAttachment
- [x] Status workflow
- [x] Claim use case
- [x] Claim HTTP handler

### Phase 5: Pricing Engine (Grule)
- [ ] Rule file structure
- [ ] Coverage rules
- [ ] Franchise rules
- [ ] Deduction rules

### Phase 6: Provider Settlement ✅
- [x] Center entity (center.go)
- [x] CenterContract
- [x] Package entity
- [x] Settlement entity
- [ ] Payment workflow APIs

### Phase 7: Integrations
- [ ] Tamin API client
- [ ] Sepas API client
- [ ] IRC integration

### Phase 8: Frontend (Vue.js) ✅
- [x] Vue 3 project setup
- [x] Sunset Theme from welfare-V2
- [x] Login page
- [x] Dashboard
- [x] Claims management
- [x] Packages view
- [x] Centers view
- [x] Settlements view
- [x] Members inquiry
- [x] Reports view
- [x] Users management
- [x] Settings view

### Phase 9: Deployment ✅
- [x] Docker Compose
- [x] Dockerfile (backend)
- [x] Dockerfile (frontend)
- [x] Nginx config for server
- [x] Deploy script

## Deployment

### Server Info
- IP: 37.152.174.87
- Domain: ria.jafamhis.ir
- Path: /tpa
- Port: 8086 (via nginx)

### Existing Routes (DO NOT MODIFY)
- / (root)
- /commission
- /hotel
- /welfare

### Deploy Commands
```bash
# On server
ssh root@37.152.174.87
cd /root/projects/tpa
git pull
docker-compose up -d --build

# Nginx config (add to existing)
location /tpa {
    proxy_pass http://localhost:8086/tpa;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
}

# Reload nginx
sudo nginx -t && sudo systemctl reload nginx
```

## Styling

Frontend uses Sunset Theme from welfare-V2:
- Primary: Coral Red (#ff6b6b)
- Secondary: Turquoise (#00d2d3)
- Accent: Golden Yellow (#feca57)
- Purple: Vibrant Purple (#a55eea)
- Font: Vazirmatn (Persian)
- RTL Support

## Contact

Project maintained for Bank Melli Iran welfare department.
