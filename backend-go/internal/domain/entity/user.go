package entity

import (
	"time"
)

// User - کاربر سیستم (Simplified)
type User struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_users_tenant" json:"tenant_id"`

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

	// وضعیت
	IsActive         bool       `gorm:"default:true" json:"is_active"`
	IsEmailVerified  bool       `gorm:"default:false" json:"is_email_verified"`
	IsMobileVerified bool       `gorm:"default:false" json:"is_mobile_verified"`
	LastLoginAt      *time.Time `json:"last_login_at"`

	// Relations
	Tenant *Insurer `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// Role - نقش کاربری
type Role struct {
	BaseModel

	Name        RoleName    `gorm:"size:50;uniqueIndex" json:"name"`
	TitleFa     string      `gorm:"size:100" json:"title_fa"`
	Level       int         `gorm:"default:0" json:"level"`
	IsSystem    bool        `gorm:"default:false" json:"is_system"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission - مجوز دسترسی
type Permission struct {
	BaseModel

	Name     PermissionName `gorm:"size:100;uniqueIndex" json:"name"`
	TitleFa  string         `gorm:"size:100" json:"title_fa"`
	Module   string         `gorm:"size:50" json:"module"`
	IsActive bool           `gorm:"default:true" json:"is_active"`
}

// UserRefreshToken - توکن‌های تازه‌سازی کاربران
type UserRefreshToken struct {
	BaseModel

	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Token     string     `gorm:"size:500;uniqueIndex" json:"token"`
	ExpiresAt time.Time  `json:"expires_at"`
	IsRevoked bool       `gorm:"default:false" json:"is_revoked"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName implementations
func (User) TableName() string { return "users" }
func (Role) TableName() string { return "roles" }
func (Permission) TableName() string { return "permissions" }
func (UserRefreshToken) TableName() string { return "user_refresh_tokens" }
