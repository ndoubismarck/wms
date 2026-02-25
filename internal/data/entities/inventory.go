package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID                 string          `gorm:"type:varchar(150);primaryKey" json:"id"`
	BinID              string          `gorm:"type:varchar(150);not null;uniqueIndex:idx_inventory_variant_bin" json:"bin_id"`
	VariantID          string          `gorm:"type:varchar(150);not null;uniqueIndex:idx_inventory_variant_bin" json:"variant_id"`
	QuantityOnHand     uint64          `gorm:"not null;default:0" json:"quantity_on_hand"`
	QuantityDamaged    uint64          `gorm:"not null;default:0" json:"quantity_damaged"`
	QuantityReserved   uint64          `gorm:"not null;default:0" json:"quantity_reserved"`
	QuantityIncoming   uint64          `gorm:"not null;default:0" json:"quantity_incoming"`
	QuantityOutgoing   uint64          `gorm:"not null;default:0" json:"quantity_outgoing"`
	QuantityAvailable  uint64          `gorm:"-" json:"quantity_available"`
	QuantityForecasted uint64          `gorm:"-" json:"quantity_forecasted"`
	Location           string          `gorm:"->" json:"location"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	Bin                *LocationBin    `gorm:"foreignKey:BinID;references:ID" json:"bin,omitempty"`
	Variant            *ProductVariant `gorm:"foreignKey:VariantID;references:ID" json:"variant,omitempty"`
}

type InventorySummary struct {
	LocationID       string `json:"location_id"`
	VariantID        string `json:"variant_id"`
	QuantityOnHand   uint64 `json:"quantity_on_hand"`
	QuantityReserved uint64 `json:"quantity_reserved"`
	QuantityIncoming uint64 `json:"quantity_incoming"`
	QuantityOutgoing uint64 `json:"quantity_outgoing"`
}

func (e *Inventory) TableName() string {
	return "inventory"
}

func (e *Inventory) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}

func (e *Inventory) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

func (e *Inventory) AfterFind(tx *gorm.DB) error {
	if e.QuantityOnHand > e.QuantityReserved {
		e.QuantityAvailable = e.QuantityOnHand - e.QuantityReserved
	} else {
		e.QuantityAvailable = 0
	}
	base := e.QuantityOnHand + e.QuantityIncoming
	if base > e.QuantityOutgoing {
		e.QuantityForecasted = base - e.QuantityOutgoing
	} else {
		e.QuantityForecasted = 0
	}
	return nil
}
