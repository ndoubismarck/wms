package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ShipmentStatusEnum string

const (
	ShipmentStatusDraft     ShipmentStatusEnum = "draft"
	ShipmentStatusReady     ShipmentStatusEnum = "ready"
	ShipmentStatusInTransit ShipmentStatusEnum = "in_transit"
	ShipmentStatusDelivered ShipmentStatusEnum = "delivered"
	ShipmentStatusException ShipmentStatusEnum = "exception"
)

type Shipment struct {
	ID             string             `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID     string             `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_shipments_location_number" json:"location_id"`
	ShipmentNumber string             `gorm:"type:varchar(100);not null;uniqueIndex:idx_shipments_location_number" json:"shipment_number"`
	OrderNumber    string             `gorm:"type:varchar(100);not null" json:"order_number"`
	Carrier        string             `gorm:"type:varchar(120);not null" json:"carrier"`
	Service        string             `gorm:"type:varchar(120);not null" json:"service"`
	Status         ShipmentStatusEnum `gorm:"type:varchar(25);not null;default:'draft';check(status IN ('draft','ready','in_transit','delivered','exception'))" json:"status"`
	Packages       uint64             `gorm:"not null;default:0" json:"packages"`
	TrackingCode   string             `gorm:"type:varchar(200);not null" json:"tracking_code"`
	ETA            *time.Time         `gorm:"null" json:"eta,omitempty"`
	Destination    string             `gorm:"type:varchar(200);not null" json:"destination"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	Location       *Location          `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}

func (e *Shipment) TableName() string {
	return "shipments"
}

func (e *Shipment) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
