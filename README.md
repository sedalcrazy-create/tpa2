# TPA System - سامانه مدیریت اسناد درمانی

سامانه جامع مدیریت اسناد درمانی (Third Party Administrator) برای بیمه‌های تکمیلی درمان

## ویژگی‌ها

- مدیریت ادعاهای درمانی (12 نوع ادعا)
- مدیریت بسته‌های اسناد از مراکز درمانی
- تسویه حساب با مراکز طرف قرارداد
- استعلام بیمه‌شدگان
- گزارشات تحلیلی
- پشتیبانی از چند بیمه‌گر (Multi-tenant)

## معماری

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Vue.js 3      │────▶│   Go + Fiber    │────▶│  PostgreSQL     │
│   Frontend      │     │   Backend API   │     │  Database       │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               │
                               ▼
                        ┌─────────────────┐
                        │   Grule Rule    │
                        │   Engine        │
                        └─────────────────┘
```

## تکنولوژی‌ها

### Backend
- Go 1.22+
- Fiber (Web Framework)
- GORM (ORM)
- Grule Rule Engine (قوانین قیمت‌گذاری)
- PostgreSQL
- Redis

### Frontend
- Vue.js 3
- TypeScript
- Pinia (State Management)
- Vue Router
- SCSS (Sunset Theme)

## راه‌اندازی محلی

### Backend
```bash
cd backend-go
cp .env.example .env
# Edit .env with your settings
go mod download
go run cmd/api/main.go
```

### Frontend
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## دیپلوی با Docker

```bash
# Build and run all services
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop services
docker-compose down
```

## ساختار پروژه

```
TPA2/
├── backend-go/
│   ├── cmd/api/                 # Entry point
│   ├── internal/
│   │   ├── config/              # Configuration
│   │   ├── domain/
│   │   │   ├── entity/          # Domain entities
│   │   │   └── repository/      # Repository interfaces
│   │   ├── usecase/             # Business logic
│   │   ├── infrastructure/      # Database, external services
│   │   └── delivery/http/       # HTTP handlers
│   ├── Dockerfile
│   └── docker-compose.yml
├── frontend/
│   ├── src/
│   │   ├── assets/              # Styles, images
│   │   ├── components/          # Vue components
│   │   ├── layouts/             # Layout components
│   │   ├── views/               # Page views
│   │   ├── stores/              # Pinia stores
│   │   ├── services/            # API services
│   │   └── router/              # Vue Router
│   └── Dockerfile
├── docker-compose.yml           # Main compose file
├── nginx-server.conf            # Nginx config for server
└── deploy.sh                    # Deployment script
```

## انواع ادعا

| کد | نوع |
|----|-----|
| 1 | داروخانه |
| 2 | بستری |
| 3 | دندانپزشکی |
| 4 | ویزیت |
| 5 | آزمایشگاه |
| 6 | تصویربرداری |
| 7 | فیزیوتراپی |
| 8 | جراحی سرپایی |
| 9 | اورژانس |
| 10 | تجهیزات پزشکی |
| 13 | تزریقات |
| 15 | سرپایی بیمارستان |

## وضعیت ادعا

```
ثبت ← منتظر ارزیابی ← منتظر تایید ← ارسال به مالی ← آرشیو
                ↓
            عودت شده
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/login | ورود |
| GET | /api/v1/claims | لیست ادعاها |
| POST | /api/v1/claims | ایجاد ادعا |
| GET | /api/v1/claims/:id | جزئیات ادعا |
| POST | /api/v1/claims/:id/examine | ارزیابی |
| POST | /api/v1/claims/:id/approve | تایید |
| GET | /api/v1/packages | لیست بسته‌ها |
| GET | /api/v1/centers | لیست مراکز |
| GET | /api/v1/settlements | تسویه حساب‌ها |
| GET | /api/v1/members/inquiry | استعلام بیمه‌شده |

## تنظیمات Nginx

پس از دیپلوی، این location را به nginx config سرور اضافه کنید:

```nginx
location /tpa {
    proxy_pass http://localhost:8086/tpa;
    proxy_http_version 1.1;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

## نقش‌های کاربری

| نقش | توضیح |
|-----|-------|
| system_admin | مدیر سیستم |
| insurer_admin | مدیر بیمه‌گر |
| supervisor | سرپرست |
| claim_examiner | ارزیاب |
| drug_examiner | ارزیاب دارو |
| financial_officer | کارشناس مالی |
| center_user | کاربر مرکز |
| report_viewer | مشاهده‌کننده گزارش |

## مجوز

این پروژه تحت مجوز اختصاصی است.

---

طراحی و توسعه: 2026
