package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type UserRoleEnum string

const (
	UserRoleNone       UserRoleEnum = "none"
	UserRoleAdmin                   = "admin"
	UserRoleEditor                  = "editor"
	UserRoleSubscriber              = "subscriber"
)

type UserStatusEnum string

const (
	UserStatusActive    UserStatusEnum = "active"
	UserStatusInactive  UserStatusEnum = "inactive"
	UserStatusSuspended UserStatusEnum = "suspended"
)

type User struct {
	ID                 string         `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationD          *string        `gorm:"type:varchar(150);null;index" json:"location_id,omitempty"`
	Role               UserRoleEnum   `gorm:"type:varchar(150);check(role IN ('none', 'admin', 'editor', 'subscriber'))" json:"role"`
	Status             UserStatusEnum `gorm:"type:varchar(150);check(status IN ('active', 'inactive', 'suspended'))" json:"status"`
	FirstName          string         `gorm:"type:varchar(100)" json:"first_name"`
	LastName           string         `gorm:"type:varchar(100)" json:"last_name"`
	PhoneNumber        string         `gorm:"type:varchar(25)" json:"phone_number"`
	EmailAddress       string         `gorm:"type:varchar(150);unique" json:"email_address"`
	Password           string         `gorm:"type:varchar(200)" json:"-"`
	LastLoginIPAddress *string        `gorm:"type:varchar(25)" json:"last_login_ip_address"`
	LastLoginAt        *time.Time     `json:"last_login_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (e *User) TableName() string {
	return "users"
}

func (e *User) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
