package entity

import (
	"time"
)

// User - کاربر سیستم
type User struct {
	TenantModel

	// اطلاعات ورود
	Username     string `gorm:"size:50;uniqueIndex" json:"username"`
	Email        string `gorm:"size:100;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"size:255" json:"-"`
	Mobile       string `gorm:"size:15" json:"mobile"`

	// اطلاعات شخصی
	FirstName    string `gorm:"size:100" json:"first_name"`
	LastName     string `gorm:"size:100" json:"last_name"`
	NationalCode string `gorm:"size:10" json:"national_code"`
	Avatar       string `gorm:"size:255" json:"avatar"`

	// نقش و دسترسی
	RoleID uint  `json:"role_id"`
	Role   *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`

	// مرکز درمانی (برای کاربران مرکز)
	CenterID *uint   `json:"center_id"`
	Center   *Center `gorm:"foreignKey:CenterID" json:"center,omitempty"`

	// واحد کاری
	WorkUnitID *uint     `json:"work_unit_id"`
	WorkUnit   *WorkUnit `gorm:"foreignKey:WorkUnitID" json:"work_unit,omitempty"`

	// استان (برای محدودیت دسترسی جغرافیایی)
	ProvinceID *uint     `json:"province_id"`
	Province   *Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`

	// وضعیت
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IsEmailVerified  bool       `gorm:"default:false" json:"is_email_verified"`
	IsMobileVerified bool       `gorm:"default:false" json:"is_mobile_verified"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`

	// تنظیمات
	MustChangePassword bool   `gorm:"default:false" json:"must_change_password"`
	FailedLoginCount   int    `gorm:"default:0" json:"failed_login_count"`
	LockedUntil        *time.Time `json:"locked_until"`

	// Relations
	RefreshTokens []UserRefreshToken `json:"-"`
	ActivityLogs  []UserActivityLog  `json:"-"`
}

// Role - نقش کاربری
type Role struct {
	BaseModel

	Name        string `gorm:"size:50;uniqueIndex" json:"name"`
	TitleFa     string `gorm:"size:100" json:"title_fa"`
	Description string `gorm:"size:500" json:"description"`

	// سطح دسترسی
	Level    int  `json:"level"`     // سطح (برای سلسله مراتب)
	IsSystem bool `gorm:"default:false" json:"is_system"` // نقش سیستمی

	IsActive bool `gorm:"default:true" json:"is_active"`

	// Relations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission - دسترسی
type Permission struct {
	BaseModel

	Name        string `gorm:"size:100;uniqueIndex" json:"name"` // e.g., "claims.create"
	TitleFa     string `gorm:"size:100" json:"title_fa"`
	Description string `gorm:"size:500" json:"description"`
	Module      string `gorm:"size:50" json:"module"` // ماژول

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// UserRefreshToken - توکن رفرش
type UserRefreshToken struct {
	BaseModel

	UserID uint `gorm:"index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	Token     string    `gorm:"size:500;uniqueIndex" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`

	// اطلاعات دستگاه
	UserAgent string `gorm:"size:500" json:"user_agent"`
	IPAddress string `gorm:"size:45" json:"ip_address"`
}

// UserActivityLog - لاگ فعالیت کاربر
type UserActivityLog struct {
	BaseModel

	UserID uint `gorm:"index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	Action      string `gorm:"size:50" json:"action"`       // login, logout, create, update, delete
	Module      string `gorm:"size:50" json:"module"`       // claims, users, centers, ...
	EntityType  string `gorm:"size:50" json:"entity_type"`  // نوع موجودیت
	EntityID    *uint  `json:"entity_id"`                   // شناسه موجودیت
	Description string `gorm:"size:1000" json:"description"`
	OldValue    string `gorm:"type:text" json:"old_value"`  // مقدار قبلی (JSON)
	NewValue    string `gorm:"type:text" json:"new_value"`  // مقدار جدید (JSON)

	// اطلاعات درخواست
	IPAddress string `gorm:"size:45" json:"ip_address"`
	UserAgent string `gorm:"size:500" json:"user_agent"`
	RequestID string `gorm:"size:50" json:"request_id"`
}

// AuditLog - لاگ ممیزی (سطح سیستم)
type AuditLog struct {
	BaseModel
	TenantID uint `gorm:"index" json:"tenant_id"`

	UserID      *uint  `json:"user_id"`
	Action      string `gorm:"size:50" json:"action"`
	Module      string `gorm:"size:50" json:"module"`
	EntityType  string `gorm:"size:50" json:"entity_type"`
	EntityID    *uint  `json:"entity_id"`
	Description string `gorm:"size:1000" json:"description"`
	OldValue    string `gorm:"type:text" json:"old_value"`
	NewValue    string `gorm:"type:text" json:"new_value"`
	IPAddress   string `gorm:"size:45" json:"ip_address"`
	UserAgent   string `gorm:"size:500" json:"user_agent"`
	RequestID   string `gorm:"size:50" json:"request_id"`
	Severity    uint8  `json:"severity"` // 1: info, 2: warning, 3: error, 4: critical
}

// Notification - اعلان
type Notification struct {
	TenantModel

	// نوع
	Type     string `gorm:"size:50" json:"type"`      // info, warning, error, success
	Category string `gorm:"size:50" json:"category"`  // claim, package, payment, system
	Priority uint8  `json:"priority"`                 // 1: low, 2: medium, 3: high, 4: urgent

	// محتوا
	Title   string `gorm:"size:200" json:"title"`
	Message string `gorm:"size:2000" json:"message"`
	Link    string `gorm:"size:500" json:"link"` // لینک مرتبط

	// مرتبط با
	EntityType string `gorm:"size:50" json:"entity_type"`
	EntityID   *uint  `json:"entity_id"`

	// زمان‌بندی
	ScheduledAt *time.Time `json:"scheduled_at"`
	ExpiresAt   *time.Time `json:"expires_at"`

	// Relations
	Recipients []NotificationRecipient `json:"recipients,omitempty"`
}

// NotificationRecipient - گیرنده اعلان
type NotificationRecipient struct {
	BaseModel

	NotificationID uint         `gorm:"index" json:"notification_id"`
	Notification   Notification `gorm:"foreignKey:NotificationID" json:"-"`

	UserID uint `gorm:"index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"-"`

	// وضعیت
	IsRead   bool       `gorm:"default:false" json:"is_read"`
	ReadAt   *time.Time `json:"read_at"`
	IsSent   bool       `gorm:"default:false" json:"is_sent"`
	SentAt   *time.Time `json:"sent_at"`
	Channel  string     `gorm:"size:20" json:"channel"` // app, email, sms
}

// SystemSetting - تنظیمات سیستم
type SystemSetting struct {
	BaseModel
	TenantID *uint `gorm:"index" json:"tenant_id"` // null = global

	Key         string `gorm:"size:100;uniqueIndex:idx_tenant_key" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Type        string `gorm:"size:20" json:"type"` // string, int, bool, json
	Description string `gorm:"size:500" json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"` // قابل دسترس در فرانت
}

// Predefined roles
const (
	RoleSystemAdmin       = "system_admin"
	RoleInsurerAdmin      = "insurer_admin"
	RoleSupervisor        = "supervisor"
	RoleClaimExaminer     = "claim_examiner"
	RoleDrugExaminer      = "drug_examiner"
	RoleFinancialOfficer  = "financial_officer"
	RoleCenterUser        = "center_user"
	RoleReportViewer      = "report_viewer"
)

// Predefined permissions
const (
	// Claims
	PermClaimCreate          = "claims.create"
	PermClaimRead            = "claims.read"
	PermClaimUpdate          = "claims.update"
	PermClaimDelete          = "claims.delete"
	PermClaimExamine         = "claims.examine"
	PermClaimApprove         = "claims.approve"
	PermClaimReject          = "claims.reject"

	// Packages
	PermPackageCreate        = "packages.create"
	PermPackageRead          = "packages.read"
	PermPackageUpdate        = "packages.update"
	PermPackageDelete        = "packages.delete"
	PermPackageExamine       = "packages.examine"
	PermPackageApprove       = "packages.approve"

	// Centers
	PermCenterCreate         = "centers.create"
	PermCenterRead           = "centers.read"
	PermCenterUpdate         = "centers.update"
	PermCenterDelete         = "centers.delete"

	// Settlements
	PermSettlementCreate     = "settlements.create"
	PermSettlementRead       = "settlements.read"
	PermSettlementApprove    = "settlements.approve"

	// Users
	PermUserCreate           = "users.create"
	PermUserRead             = "users.read"
	PermUserUpdate           = "users.update"
	PermUserDelete           = "users.delete"

	// Reports
	PermReportView           = "reports.view"
	PermReportExport         = "reports.export"

	// Settings
	PermSettingsRead         = "settings.read"
	PermSettingsUpdate       = "settings.update"
)
