package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type SupplierTypeEnum string

type SupplierStatusEnum string

const (
	SupplierTypeRawMaterial  SupplierTypeEnum = "raw_material"
	SupplierTypePackaging    SupplierTypeEnum = "packaging"
	SupplierTypeFinishedGood SupplierTypeEnum = "finished_goods"
	SupplierTypeService      SupplierTypeEnum = "service"
)

const (
	SupplierStatusApproved  SupplierStatusEnum = "approved"
	SupplierStatusProbation SupplierStatusEnum = "probation"
	SupplierStatusBlocked   SupplierStatusEnum = "blocked"
)

type Supplier struct {
	ID                 string             `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID         string             `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_suppliers_location_name" json:"location_id"`
	Name               string             `gorm:"type:varchar(200);not null;uniqueIndex:idx_suppliers_location_name" json:"name"`
	SupplierType       SupplierTypeEnum   `gorm:"type:varchar(30);not null;default:'service';check(supplier_type IN ('raw_material','packaging','finished_goods','service'))" json:"supplier_type"`
	Status             SupplierStatusEnum `gorm:"type:varchar(25);not null;default:'approved';check(status IN ('approved','probation','blocked'))" json:"status"`
	ContactName        string             `gorm:"type:varchar(150);not null" json:"contact_name"`
	Email              string             `gorm:"type:varchar(200);not null" json:"email"`
	LeadTimeDays       uint64             `gorm:"not null;default:0" json:"lead_time_days"`
	OnTimeRate         float64            `gorm:"type:double;not null;default:0" json:"on_time_rate"`
	OpenPurchaseOrders uint64             `gorm:"not null;default:0" json:"open_purchase_orders"`
	City               string             `gorm:"type:varchar(120);not null" json:"city"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	Location           *Location          `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}

func (e *Supplier) TableName() string {
	return "suppliers"
}

func (e *Supplier) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
