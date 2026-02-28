package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type CustomerTierEnum string

type CustomerStatusEnum string

const (
	CustomerTierStandard   CustomerTierEnum = "standard"
	CustomerTierGrowth     CustomerTierEnum = "growth"
	CustomerTierEnterprise CustomerTierEnum = "enterprise"
)

const (
	CustomerStatusActive    CustomerStatusEnum = "active"
	CustomerStatusPaused    CustomerStatusEnum = "paused"
	CustomerStatusChurnRisk CustomerStatusEnum = "churn_risk"
)

type Customer struct {
	ID            string             `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID    string             `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_customers_location_email" json:"location_id"`
	Name          string             `gorm:"type:varchar(200);not null" json:"name"`
	Email         string             `gorm:"type:varchar(200);not null;uniqueIndex:idx_customers_location_email" json:"email"`
	Phone         string             `gorm:"type:varchar(50);not null" json:"phone"`
	City          string             `gorm:"type:varchar(120);not null" json:"city"`
	Tier          CustomerTierEnum   `gorm:"type:varchar(25);not null;default:'standard';check(tier IN ('standard','growth','enterprise'))" json:"tier"`
	Status        CustomerStatusEnum `gorm:"type:varchar(25);not null;default:'active';check(status IN ('active','paused','churn_risk'))" json:"status"`
	TotalOrders   uint64             `gorm:"not null;default:0" json:"total_orders"`
	LifetimeValue float64            `gorm:"type:double;not null;default:0" json:"lifetime_value"`
	LastOrderAt   *time.Time         `gorm:"null" json:"last_order_at,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Location      *Location          `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}

func (e *Customer) TableName() string {
	return "customers"
}

func (e *Customer) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
