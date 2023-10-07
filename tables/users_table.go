package tables

import (
	"time"

	"gorm.io/gorm"
)

// ID            int            `gorm:"primary_key:auto_increment" json:"id"`
// Name          string         `gorm:"type:varchar(255)" json:"name"`
// Phone         string         `gorm:"type:varchar(255)" json:"phone"`
// Email         string         `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
// Password      string         `gorm:"->;<-;not null" json:"-"`
// Avatar        string         `gorm:"type:varchar(255)" json:"avatar"`
// RememberToken string         `gorm:"type:varchar(255)" json:"remember_token"`
// CreatedAt     time.Time      `json:"created_at,omitemty"`
// UpdatedAt     time.Time      `json:"updated_at,omitemty"`
// DeletedAt     gorm.DeletedAt `gorm:"index"`

type Users struct {
	ID        int
	Name      string
	Username  string `json:"username" `
	Phone     string `gorm:"type:varchar(13);default:null" json:"phone"`
	Email     string `gorm:"type:varchar(255);default:null" json:"email"`
	Password  string `gorm:"default:null"`
	RoleID    int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserData struct {
	ID              int
	Name            string
	Phone           string `gorm:"type:varchar(13);default:null" json:"phone"`
	Email           string `gorm:"type:varchar(255);default:null" json:"email"`
	Password        string
	Avatar          string
	DateOfBirth     string `sql:"default:null" gorm:"default:null" `
	GenderID        int
	GenderName      string
	RememberToken   string
	IsEmailVerified bool `gorm:"default:false"`
	isPhoneVerified bool `gorm:"default:false"`
	Status          bool `gorm:"default:true"`
	EncryptCode     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	RoleID          int
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type RoleUser struct {
	ID     int
	RoleID int
	UserID int
}

type UserMember struct {
	ID              int
	Name            string
	Phone           string
	Email           string
	RememberToken   string
	Password        string
	IsEmailVerified bool
	IsPhoneVerified bool
	GenderID        int
	DateOfBirth     string
	Avatar          string
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type Genders struct {
	ID          int
	Code        string
	Translation string
}

type UserInvites struct {
	UserID         int
	Phone          string
	OrganizationID int
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type UserInvitesOrg struct {
	ID             int
	UserID         int
	Phone          string
	OrganizationID int
}

type UserOTP struct {
	ID        int
	AppTypeID int
	UserID    int
	Phone     string
	OtpCode   string
	Status    bool
	RoleID    int
}
