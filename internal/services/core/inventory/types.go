package inventory

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"time"
)

type (
	AddData struct {
		LocationID       string `json:"-" validate:"required"`
		BinID            string `json:"bin_id" validate:"required"`
		VariantID        string `json:"variant_id" validate:"required"`
		QuantityOnHand   uint64 `json:"quantity_on_hand"`
		QuantityDamaged  uint64 `json:"quantity_damaged"`
		QuantityReserved uint64 `json:"quantity_reserved"`
		QuantityIncoming uint64 `json:"quantity_incoming"`
		QuantityOutgoing uint64 `json:"quantity_outgoing"`
	}
	AddResult struct {
		Code       types.ServiceResultCode
		Payload    AddResultPayload
		Validation types.ValidationResult
	}
	AddResultPayload struct {
		Inventory entities.Inventory `json:"inventory"`
	}
)

type (
	UpdateData struct {
		LocationID       string `json:"-" validate:"required"`
		ID               string `json:"id" validate:"required"`
		BinID            string `json:"bin_id"`
		VariantID        string `json:"variant_id"`
		QuantityOnHand   uint64 `json:"quantity_on_hand"`
		QuantityDamaged  uint64 `json:"quantity_damaged"`
		QuantityReserved uint64 `json:"quantity_reserved"`
		QuantityIncoming uint64 `json:"quantity_incoming"`
		QuantityOutgoing uint64 `json:"quantity_outgoing"`
	}
	UpdateResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateResultPayload
		Validation types.ValidationResult
	}
	UpdateResultPayload struct {
		Inventory entities.Inventory `json:"inventory"`
	}
)

type (
	DeleteData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetOneData struct {
		LocationID    string     `json:"-" form:"-" validate:"required"`
		ID            *string    `json:"id" form:"id"`
		BinID         *string    `json:"bin_id" form:"bin_id"`
		Related       bool       `json:"related" form:"related"`
		ToUpdatedAt   *time.Time `json:"to_updated_at" form:"to_updated_at"`
		FromUpdatedAt *time.Time `json:"from_updated_at" form:"from_updated_at"`
	}
	GetOneResult struct {
		Code    types.ServiceResultCode
		Payload GetOneResultPayload
	}
	GetOneResultPayload struct {
		Inventory entities.Inventory `json:"inventory"`
	}
)

type (
	GetManyData struct {
		types.PaginationParams
		LocationID    string     `json:"-" form:"-" validate:"required"`
		BinID         *string    `json:"bin_id" form:"bin_id"`
		Related       bool       `json:"related" form:"related"`
		ToUpdatedAt   *time.Time `json:"to_updated_at" form:"to_updated_at"`
		FromUpdatedAt *time.Time `json:"from_updated_at" form:"from_updated_at"`
	}
	GetManyResult struct {
		Code       types.ServiceResultCode
		Payload    GetManyResultPayload
		Pagination *types.PaginationResult
	}
	GetManyResultPayload struct {
		Inventory []entities.Inventory `json:"inventory"`
	}
)

type (
	GetSummaryData struct {
		LocationID    string     `json:"-" form:"-" validate:"required"`
		ID            *string    `json:"id" form:"id"`
		BinID         *string    `json:"bin_id" form:"bin_id"`
		ToUpdatedAt   *time.Time `json:"to_updated_at" form:"to_updated_at"`
		FromUpdatedAt *time.Time `json:"from_updated_at" form:"from_updated_at"`
	}
	GetSummaryResult struct {
		Code    types.ServiceResultCode
		Payload GetSummaryResultPayload
	}
	GetSummaryResultPayload struct {
		Summary []entities.InventorySummary `json:"summary"`
	}
)

type (
	AddMovementData struct {
		LocationID   string                             `json:"-" validate:"required"`
		InventoryID  string                             `json:"inventory_id" validate:"required"`
		ReferenceID  string                             `json:"reference_id" validate:"required"`
		MovementType entities.InventoryMovementTypeEnum `json:"movement_type" validate:"required"`
		Quantity     uint64                             `json:"quantity" validate:"required"`
		Reason       *string                            `json:"reason,omitempty"`
	}
	AddMovementResult struct {
		Code       types.ServiceResultCode
		Payload    AddMovementResultPayload
		Validation types.ValidationResult
	}
	AddMovementResultPayload struct {
		Movement entities.InventoryMovement `json:"movement"`
	}
)

type (
	UpdateMovementData struct {
		LocationID   string                             `json:"-" validate:"required"`
		ID           string                             `json:"id" validate:"required"`
		InventoryID  string                             `json:"inventory_id"`
		ReferenceID  string                             `json:"reference_id"`
		MovementType entities.InventoryMovementTypeEnum `json:"movement_type"`
		Quantity     uint64                             `json:"quantity"`
		Reason       *string                            `json:"reason,omitempty"`
	}
	UpdateMovementResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateMovementResultPayload
		Validation types.ValidationResult
	}
	UpdateMovementResultPayload struct {
		Movement entities.InventoryMovement `json:"movement"`
	}
)

type (
	DeleteMovementData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteMovementResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetMovementData struct {
		LocationID   string                               `json:"-" form:"-" validate:"required"`
		ID           *string                              `json:"id" form:"id"`
		InventoryID  *string                              `json:"inventory_id" form:"inventory_id"`
		ReferenceID  *string                              `json:"reference_id" form:"reference_id"`
		MovementType []entities.InventoryMovementTypeEnum `json:"movement_type" form:"movement_type"`
	}
	GetMovementResult struct {
		Code    types.ServiceResultCode
		Payload GetMovementResultPayload
	}
	GetMovementResultPayload struct {
		Movement entities.InventoryMovement `json:"movement"`
	}
)

type (
	GetMovementsData struct {
		types.PaginationParams
		LocationID   string                               `json:"-" form:"-" validate:"required"`
		InventoryID  *string                              `json:"inventory_id" form:"inventory_id"`
		ReferenceID  *string                              `json:"reference_id" form:"reference_id"`
		MovementType []entities.InventoryMovementTypeEnum `json:"movement_type" form:"movement_type"`
	}
	GetMovementsResult struct {
		Code       types.ServiceResultCode
		Payload    GetMovementsResultPayload
		Pagination *types.PaginationResult
	}
	GetMovementsResultPayload struct {
		Movements []entities.InventoryMovement `json:"movements"`
	}
)
