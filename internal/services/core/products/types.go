package products

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
)

type (
	attribute struct {
		String  *string                              `json:"string,omitempty"`
		Int     *int64                               `json:"int,omitempty"`
		Float   *float64                             `json:"float,omitempty"`
		Bool    *bool                                `json:"bool,omitempty"`
		Date    *string                              `json:"date,omitempty"`
		JSON    any                                  `json:"json,omitempty"`
		Range   *entities.ProductAttributeRangeValue `json:"range,omitempty"`
		Option  *string                              `json:"option,omitempty"`
		Options []string                             `json:"options,omitempty"`
	}
)

type (
	AddData struct {
		BinID         string               `json:"bin_id"`
		BrandID       string               `json:"brand_id"`
		LocationID    string               `json:"-" validate:"required"`
		SubcategoryID string               `json:"subcategory_id"`
		Name          string               `json:"name"`
		Description   string               `json:"description"`
		Images        []string             `json:"images"`
		Attributes    map[string]attribute `json:"attributes"`
	}
	AddResult struct {
		Code       types.ServiceResultCode
		Payload    AddResultPayload
		Validation types.ValidationResult
	}

	AddResultPayload struct {
		Product entities.Product `json:"product"`
	}
)

type (
	UpdateData struct {
		LocationID    string                      `json:"-" validate:"required"`
		ID            string                      `json:"id" validate:"required"`
		SubcategoryID *string                     `json:"subcategory_id,omitempty"`
		BrandID       *string                     `json:"brand_id,omitempty"`
		Name          *string                     `json:"name,omitempty"`
		Description   *string                     `json:"description,omitempty"`
		Status        *entities.ProductStatusEnum `json:"status,omitempty"`
	}
	UpdateResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateResultPayload
		Validation types.ValidationResult
	}
	UpdateResultPayload struct {
		Product entities.Product `json:"product"`
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
		LocationID string  `json:"-" validate:"required"`
		ID         *string `json:"id"`
		BinID      *string `json:"bin_id"`
	}
	GetOneResult struct {
		Code    types.ServiceResultCode
		Payload GetOneResultPayload
	}

	GetOneResultPayload struct {
		Product entities.Product `json:"product"`
	}
)

type (
	GetManyData struct {
		types.PaginationParams
		LocationID string  `json:"-" form:"-" validate:"required"`
		BinID      *string `json:"bin_id" form:"bin_id"`
	}
	GetManyResult struct {
		Code       types.ServiceResultCode
		Payload    GetManyResultPayload
		Pagination *types.PaginationResult
	}

	GetManyResultPayload struct {
		Products []entities.Product `json:"products"`
	}
)

type (
	AddCategoryData struct {
		LocationID string `json:"-" validate:"required"`
		BrandID    string `json:"brand_id"`
		Name       string `json:"name"`
	}
	AddCategoryResult struct {
		Code       types.ServiceResultCode
		Payload    AddCategoryResultPayload
		Validation types.ValidationResult
	}
	AddCategoryResultPayload struct {
		Category entities.ProductCategory `json:"category"`
	}
)

type (
	UpdateCategoryData struct {
		ID         string  `json:"-"`
		LocationID string  `json:"-"`
		BrandID    *string `json:"brand_id"`
		Name       *string `json:"name"`
	}
	UpdateCategoryResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateCategoryResultPayload
		Validation types.ValidationResult
	}
	UpdateCategoryResultPayload struct {
		Category entities.ProductCategory `json:"category"`
	}
)

type (
	DeleteCategoryData struct {
		ID         string `json:"id" validate:"required"`
		LocationID string `json:"-" validate:"required"`
	}
	DeleteCategoryResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetCategoryData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetCategoryResult struct {
		Code    types.ServiceResultCode
		Payload GetCategoryResultPayload
	}
	GetCategoryResultPayload struct {
		Category entities.ProductCategory `json:"category"`
	}
)

type (
	GetCategoriesData struct {
		types.PaginationParams
		BrandID    *string `json:"brand_id,omitempty" form:"brand_id"`
		LocationID string  `json:"-" form:"-" validate:"required"`
	}
	GetCategoriesResult struct {
		Code       types.ServiceResultCode
		Payload    GetCategoriesResultPayload
		Pagination *types.PaginationResult
	}
	GetCategoriesResultPayload struct {
		Categories []entities.ProductCategory `json:"categories"`
	}
)

type (
	AddSubcategoryData struct {
		LocationID string `json:"-" validate:"required"`
		CategoryID string `json:"category_id" validate:"required"`
		Name       string `json:"name"`
	}
	AddSubcategoryResult struct {
		Code       types.ServiceResultCode
		Payload    AddSubcategoryResultPayload
		Validation types.ValidationResult
	}
	AddSubcategoryResultPayload struct {
		Subcategory entities.ProductSubcategory `json:"subcategory"`
	}
)

type (
	UpdateSubcategoryData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
		CategoryID string `json:"category_id" validate:"required"`
		Name       string `json:"name"`
	}
	UpdateSubcategoryResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateSubcategoryResultPayload
		Validation types.ValidationResult
	}
	UpdateSubcategoryResultPayload struct {
		Subcategory entities.ProductSubcategory `json:"subcategory"`
	}
)

type (
	DeleteSubcategoryData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteSubcategoryResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetSubcategoryData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetSubcategoryResult struct {
		Code    types.ServiceResultCode
		Payload GetSubcategoryResultPayload
	}
	GetSubcategoryResultPayload struct {
		Subcategory entities.ProductSubcategory `json:"subcategory"`
	}
)

type (
	GetSubcategoriesData struct {
		types.PaginationParams
		LocationID string  `json:"-" form:"-" validate:"required"`
		CategoryID *string `json:"category_id" form:"category_id"`
	}
	GetSubcategoriesResult struct {
		Code       types.ServiceResultCode
		Payload    GetSubcategoriesResultPayload
		Pagination *types.PaginationResult
	}
	GetSubcategoriesResultPayload struct {
		Subcategories []entities.ProductSubcategory `json:"subcategories"`
	}
)

type (
	AddBrandData struct {
		LocationID string `json:"-" validate:"required"`
		Name       string `json:"name"`
	}
	AddBrandResult struct {
		Code       types.ServiceResultCode
		Payload    AddBrandResultPayload
		Validation types.ValidationResult
	}
	AddBrandResultPayload struct {
		Brand entities.ProductBrand `json:"brand"`
	}
)

type (
	UpdateBrandData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
		Name       string `json:"name"`
	}
	UpdateBrandResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateBrandResultPayload
		Validation types.ValidationResult
	}
	UpdateBrandResultPayload struct {
		Brand entities.ProductBrand `json:"brand"`
	}
)

type (
	DeleteBrandData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteBrandResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetBrandData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetBrandResult struct {
		Code    types.ServiceResultCode
		Payload GetBrandResultPayload
	}
	GetBrandResultPayload struct {
		Brand entities.ProductBrand `json:"brand"`
	}
)

type (
	GetBrandsData struct {
		types.PaginationParams
		LocationID string `json:"-" form:"-" validate:"required"`
	}
	GetBrandsResult struct {
		Code       types.ServiceResultCode
		Payload    GetBrandsResultPayload
		Pagination *types.PaginationResult
	}
	GetBrandsResultPayload struct {
		Brands []entities.ProductBrand `json:"brands"`
	}
)

type (
	AddAttributeData struct {
		LocationID      string                                   `json:"-" validate:"required"`
		IsSystem        bool                                     `json:"is_system"`
		Label           string                                   `json:"label"`
		Placeholder     string                                   `json:"placeholder"`
		AppliesTo       entities.ProductAttributeAppliesToEnum   `json:"applies_to"`
		DataType        entities.ProductAttributeDataTypeEnum    `json:"data_type"`
		FieldType       entities.ProductAttributeDisplayTypeEnum `json:"field_type"`
		ValidationRules []string                                 `json:"validation_rules"`
	}
	AddAttributeResult struct {
		Code       types.ServiceResultCode
		Payload    AddAttributeResultPayload
		Validation types.ValidationResult
	}
	AddAttributeResultPayload struct {
		Attribute entities.ProductAttributeDefinition `json:"attribute"`
	}
)

type (
	UpdateAttributeData struct {
		LocationID      string                                   `json:"-" validate:"required"`
		ID              string                                   `json:"id" validate:"required"`
		IsSystem        bool                                     `json:"is_system"`
		Label           string                                   `json:"label"`
		AppliesTo       entities.ProductAttributeAppliesToEnum   `json:"applies_to"`
		DataType        entities.ProductAttributeDataTypeEnum    `json:"data_type"`
		FieldType       entities.ProductAttributeDisplayTypeEnum `json:"field_type"`
		ValidationRules []string                                 `json:"validation_rules"`
	}
	UpdateAttributeResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateAttributeResultPayload
		Validation types.ValidationResult
	}
	UpdateAttributeResultPayload struct {
		Attribute entities.ProductAttributeDefinition `json:"attribute"`
	}
)

type (
	DeleteAttributeData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteAttributeResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetAttributeData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetAttributeResult struct {
		Code    types.ServiceResultCode
		Payload GetAttributeResultPayload
	}
	GetAttributeResultPayload struct {
		Attribute entities.ProductAttributeDefinition `json:"attribute"`
	}
)

type (
	GetAttributesData struct {
		types.PaginationParams
		LocationID string `json:"-" form:"-" validate:"required"`
	}
	GetAttributesResult struct {
		Code       types.ServiceResultCode
		Payload    GetAttributesResultPayload
		Pagination *types.PaginationResult
	}
	GetAttributesResultPayload struct {
		Attributes []entities.ProductAttributeDefinition `json:"attributes"`
	}
)

type (
	GetBarcodeData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"-"`
	}

	GetBarcodeResult struct {
		Code    types.ServiceResultCode
		Payload GetBarcodeResultPayload
	}

	GetBarcodeResultPayload struct {
		File GetBarcodeResultPayloadFile `json:"-"`
	}
	GetBarcodeResultPayloadFile struct {
		Path string `json:"-"`
	}
)
