package entities

import (
	"fmt"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type InventoryMovementTypeEnum string

const (
	InventoryMovementTypeReceive     InventoryMovementTypeEnum = "receive"
	InventoryMovementTypeShip        InventoryMovementTypeEnum = "ship"
	InventoryMovementTypeAdjust      InventoryMovementTypeEnum = "adjust"
	InventoryMovementTypeDamage      InventoryMovementTypeEnum = "damage"
	InventoryMovementTypeReserve     InventoryMovementTypeEnum = "reserve"
	InventoryMovementTypeRelease     InventoryMovementTypeEnum = "release"
	InventoryMovementTypeTransferOut InventoryMovementTypeEnum = "transfer_out"
	InventoryMovementTypeTransferIn  InventoryMovementTypeEnum = "transfer_in"
)

type InventoryReferenceTypeEnum string

const (
	InventoryReferenceTypePurchaseOrder InventoryReferenceTypeEnum = "purchase_order"
	InventoryReferenceTypeReturnIn      InventoryReferenceTypeEnum = "return_in"
	InventoryReferenceTypeProductionIn  InventoryReferenceTypeEnum = "production_in"
	InventoryReferenceTypePutaway       InventoryReferenceTypeEnum = "putaway"
	InventoryReferenceTypeSalesOrder    InventoryReferenceTypeEnum = "sales_order"
	InventoryReferenceTypeReturnOut     InventoryReferenceTypeEnum = "return_out"
	InventoryReferenceTypeTransfer      InventoryReferenceTypeEnum = "transfer"
	InventoryReferenceTypeCycleCount    InventoryReferenceTypeEnum = "cycle_count"
	InventoryReferenceTypeAdjustment    InventoryReferenceTypeEnum = "adjustment"
	InventoryReferenceTypeDamage        InventoryReferenceTypeEnum = "damage"
	InventoryReferenceTypeManual        InventoryReferenceTypeEnum = "manual"
	InventoryReferenceTypeSystem        InventoryReferenceTypeEnum = "system"
)

type InventoryMovement struct {
	ID           string                    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	InventoryID  string                    `gorm:"type:varchar(150);index;not null" json:"inventory_id"`
	ReferenceID  string                    `gorm:"type:varchar(150);not null" json:"reference_id"`
	MovementType InventoryMovementTypeEnum `gorm:"type:varchar(50);not null;check(type IN ('receive','ship','adjust','damage','reserve','release','transfer_out','transfer_in'))" json:"movement_type"`
	Quantity     uint64                    `gorm:"not null" json:"quantity"`
	Reason       *string                   `gorm:"type:varchar(750);null" json:"reason,omitempty"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	CreatedAt    time.Time                 `json:"created_at"`
	Inventory    *Inventory                `gorm:"foreignKey:InventoryID;references:ID" json:"inventory,omitempty"`
}

func (e *InventoryMovement) TableName() string {
	return "inventory_movements"
}

func (e *InventoryMovement) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	if e.Quantity == 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	return nil
}
