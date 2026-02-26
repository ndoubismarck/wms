package locations

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
)

type (
	GetLocationsData struct {
		types.PaginationParams
	}
	GetLocationsResult struct {
		Code       types.ServiceResultCode
		Payload    GetLocationsResultPayload
		Pagination *types.PaginationResult
	}
	GetLocationsResultPayload struct {
		Locations []entities.Location `json:"locations"`
	}
)

type (
	AddAisleData struct {
		LocationID  string `json:"-" validate:"required"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	AddAisleResult struct {
		Code       types.ServiceResultCode
		Payload    AddAisleResultPayload
		Validation types.ValidationResult
	}
	AddAisleResultPayload struct {
		Aisle entities.LocationAisle `json:"aisle"`
	}
)

type (
	UpdateAisleData struct {
		LocationID  string `json:"-" validate:"required"`
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	UpdateAisleResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateAisleResultPayload
		Validation types.ValidationResult
	}
	UpdateAisleResultPayload struct {
		Aisle entities.LocationAisle `json:"aisle"`
	}
)

type (
	DeleteAisleData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id"`
	}
	DeleteAisleResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetAisleData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetAisleResult struct {
		Code    types.ServiceResultCode
		Payload GetAisleResultPayload
	}
	GetAisleResultPayload struct {
		Aisle entities.LocationAisle `json:"aisle"`
	}
)

type (
	GetAislesData struct {
		types.PaginationParams
		LocationID string `json:"-" form:"-" validate:"required"`
	}
	GetAislesResult struct {
		Code       types.ServiceResultCode
		Payload    GetAislesResultPayload
		Pagination *types.PaginationResult
	}
	GetAislesResultPayload struct {
		Aisles []entities.LocationAisle `json:"aisles"`
	}
)

type (
	AddBayData struct {
		LocationID  string `json:"-" validate:"required"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		AisleID     string `json:"aisle_id"`
		Description string `json:"description"`
	}
	AddBayResult struct {
		Code       types.ServiceResultCode
		Payload    AddBayResultPayload
		Validation types.ValidationResult
	}
	AddBayResultPayload struct {
		Bay entities.LocationBay `json:"bay"`
	}
)

type (
	UpdateBayData struct {
		LocationID  string `json:"-" validate:"required"`
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		AisleID     string `json:"aisle_id"`
		Description string `json:"description"`
	}
	UpdateBayResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateBayResultPayload
		Validation types.ValidationResult
	}
	UpdateBayResultPayload struct {
		Bay entities.LocationBay `json:"bay"`
	}
)

type (
	DeleteBayData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id"`
	}
	DeleteBayResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetBayData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetBayResult struct {
		Code    types.ServiceResultCode
		Payload GetBayResultPayload
	}
	GetBayResultPayload struct {
		Bay entities.LocationBay `json:"bay"`
	}
)

type (
	GetBaysData struct {
		types.PaginationParams
		LocationID string  `json:"-" form:"-" validate:"required"`
		AisleID    *string `json:"aisle_id" form:"aisle_id"`
	}
	GetBaysResult struct {
		Code       types.ServiceResultCode
		Payload    GetBaysResultPayload
		Pagination *types.PaginationResult
	}
	GetBaysResultPayload struct {
		Bays []entities.LocationBay `json:"bays"`
	}
)

type (
	AddShelfData struct {
		LocationID  string `json:"-" validate:"required"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		BayID       string `json:"bay_id"`
		Description string `json:"description"`
	}
	AddShelfResult struct {
		Code       types.ServiceResultCode
		Payload    AddShelfResultPayload
		Validation types.ValidationResult
	}
	AddShelfResultPayload struct {
		Shelf entities.LocationShelf `json:"shelf"`
	}
)

type (
	UpdateShelfData struct {
		LocationID  string `json:"-" validate:"required"`
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		BayID       string `json:"bay_id"`
		Description string `json:"description"`
	}
	UpdateShelfResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateShelfResultPayload
		Validation types.ValidationResult
	}
	UpdateShelfResultPayload struct {
		Shelf entities.LocationShelf `json:"shelf"`
	}
)

type (
	DeleteShelfData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id"`
	}
	DeleteShelfResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetShelfData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetShelfResult struct {
		Code    types.ServiceResultCode
		Payload GetShelfResultPayload
	}
	GetShelfResultPayload struct {
		Shelf entities.LocationShelf `json:"shelf"`
	}
)

type (
	GetShelvesData struct {
		types.PaginationParams
		LocationID string  `json:"-" form:"-" validate:"required"`
		BayID      *string `json:"bay_id" form:"bay_id"`
	}
	GetShelvesResult struct {
		Code       types.ServiceResultCode
		Payload    GetShelvesResultPayload
		Pagination *types.PaginationResult
	}
	GetShelvesResultPayload struct {
		Shelves []entities.LocationShelf `json:"shelves"`
	}
)

type (
	AddShelfLevelData struct {
		LocationID  string `json:"-" validate:"required"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		ShelfID     string `json:"shelf_id"`
		Description string `json:"description"`
	}
	AddShelfLevelResult struct {
		Code       types.ServiceResultCode
		Payload    AddShelfLevelResultPayload
		Validation types.ValidationResult
	}
	AddShelfLevelResultPayload struct {
		ShelfLevel entities.LocationShelfLevel `json:"shelf_level"`
	}
)

type (
	UpdateShelfLevelData struct {
		LocationID  string `json:"-" validate:"required"`
		ID          string `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		ShelfID     string `json:"shelf_id"`
		Description string `json:"description"`
	}
	UpdateShelfLevelResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateShelfLevelResultPayload
		Validation types.ValidationResult
	}
	UpdateShelfLevelResultPayload struct {
		ShelfLevel entities.LocationShelfLevel `json:"shelf_level"`
	}
)

type (
	DeleteShelfLevelData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id"`
	}
	DeleteShelfLevelResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetShelfLevelData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetShelfLevelResult struct {
		Code    types.ServiceResultCode
		Payload GetShelfLevelResultPayload
	}
	GetShelfLevelResultPayload struct {
		ShelfLevel entities.LocationShelfLevel `json:"shelf_level"`
	}
)

type (
	GetShelfLevelsData struct {
		types.PaginationParams
		LocationID string  `json:"-" form:"-" validate:"required"`
		ShelfID    *string `json:"shelf_id" form:"shelf_id"`
	}
	GetShelfLevelsResult struct {
		Code       types.ServiceResultCode
		Payload    GetShelfLevelsResultPayload
		Pagination *types.PaginationResult
	}
	GetShelfLevelsResultPayload struct {
		ShelfLevels []entities.LocationShelfLevel `json:"shelf_levels"`
	}
)

type (
	AddBinData struct {
		LocationID   string `json:"-" validate:"required"`
		Code         string `json:"code"`
		Name         string `json:"name"`
		ShelfLevelID string `json:"shelf_level_id"`
		Description  string `json:"description"`
	}
	AddBinResult struct {
		Code       types.ServiceResultCode
		Payload    AddBinResultPayload
		Validation types.ValidationResult
	}
	AddBinResultPayload struct {
		Bin entities.LocationBin `json:"bin"`
	}
)

type (
	UpdateBinData struct {
		LocationID   string `json:"-" validate:"required"`
		ID           string `json:"id"`
		Code         string `json:"code"`
		Name         string `json:"name"`
		ShelfLevelID string `json:"shelf_level_id"`
		Description  string `json:"description"`
	}
	UpdateBinResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateBinResultPayload
		Validation types.ValidationResult
	}
	UpdateBinResultPayload struct {
		Bin entities.LocationBin `json:"bin"`
	}
)

type (
	DeleteBinData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id"`
	}
	DeleteBinResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetBinData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetBinResult struct {
		Code    types.ServiceResultCode
		Payload GetBinResultPayload
	}
	GetBinResultPayload struct {
		Bin entities.LocationBin `json:"bin"`
	}
)

type (
	GetBinsData struct {
		types.PaginationParams
		LocationID   string  `json:"-" form:"-" validate:"required"`
		ShelfLevelID *string `json:"shelf_level_id" form:"shelf_level_id"`
	}
	GetBinsResult struct {
		Code       types.ServiceResultCode
		Payload    GetBinsResultPayload
		Pagination *types.PaginationResult
	}
	GetBinsResultPayload struct {
		Bins []entities.LocationBin `json:"bins"`
	}
)

type (
	GetNextCodeResultPayload struct {
		Code string `json:"code"`
	}
)

type (
	GetNextAisleCodeData struct {
		LocationID string `json:"-" form:"-" validate:"required"`
	}
	GetNextAisleCodeResult struct {
		Code       types.ServiceResultCode
		Payload    GetNextCodeResultPayload
		Validation types.ValidationResult
	}
)

type (
	GetNextBayCodeData struct {
		LocationID string `json:"-" form:"-" validate:"required"`
		AisleID    string `json:"aisle_id" form:"aisle_id" validate:"required"`
	}
	GetNextBayCodeResult struct {
		Code       types.ServiceResultCode
		Payload    GetNextCodeResultPayload
		Validation types.ValidationResult
	}
)

type (
	GetNextShelfCodeData struct {
		LocationID string `json:"-" form:"-" validate:"required"`
		BayID      string `json:"bay_id" form:"bay_id" validate:"required"`
	}
	GetNextShelfCodeResult struct {
		Code       types.ServiceResultCode
		Payload    GetNextCodeResultPayload
		Validation types.ValidationResult
	}
)

type (
	GetNextShelfLevelCodeData struct {
		LocationID string `json:"-" form:"-" validate:"required"`
		ShelfID    string `json:"shelf_id" form:"shelf_id" validate:"required"`
	}
	GetNextShelfLevelCodeResult struct {
		Code       types.ServiceResultCode
		Payload    GetNextCodeResultPayload
		Validation types.ValidationResult
	}
)

type (
	GetNextBinCodeData struct {
		LocationID   string `json:"-" form:"-" validate:"required"`
		ShelfLevelID string `json:"shelf_level_id" form:"shelf_level_id" validate:"required"`
	}
	GetNextBinCodeResult struct {
		Code       types.ServiceResultCode
		Payload    GetNextCodeResultPayload
		Validation types.ValidationResult
	}
)
