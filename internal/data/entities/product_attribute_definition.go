package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductAttributeAppliesToEnum string

const (
	ProductAttributeAppliesToProduct ProductAttributeAppliesToEnum = "product"
	ProductAttributeAppliesToVariant ProductAttributeAppliesToEnum = "variant"
)

type ProductAttributeDataTypeEnum string

const (
	ProductAttributeTypeInt     ProductAttributeDataTypeEnum = "int"
	ProductAttributeTypeFloat   ProductAttributeDataTypeEnum = "float"
	ProductAttributeTypeBool    ProductAttributeDataTypeEnum = "bool"
	ProductAttributeTypeDate    ProductAttributeDataTypeEnum = "date"
	ProductAttributeTypeJSON    ProductAttributeDataTypeEnum = "json"
	ProductAttributeTypeOption  ProductAttributeDataTypeEnum = "option"
	ProductAttributeTypeString  ProductAttributeDataTypeEnum = "string"
	ProductAttributeTypeOptions ProductAttributeDataTypeEnum = "options"
)

type ProductAttributeDisplayTypeEnum string

const (
	DisplayText       ProductAttributeDisplayTypeEnum = "text"
	DisplayTextarea   ProductAttributeDisplayTypeEnum = "textarea"
	DisplayNumber     ProductAttributeDisplayTypeEnum = "number"
	DisplayCheckbox   ProductAttributeDisplayTypeEnum = "checkbox"
	DisplayRadio      ProductAttributeDisplayTypeEnum = "radio"
	DisplaySelect     ProductAttributeDisplayTypeEnum = "select"
	DisplayMultiCheck ProductAttributeDisplayTypeEnum = "multi_check"
	DisplayDate       ProductAttributeDisplayTypeEnum = "date"
)

type ProductAttributeDefinition struct {
	ID              string                          `gorm:"type:varchar(150);primaryKey" json:"id"`
	LocationD       string                          `gorm:"type:varchar(150);index" json:"location_id"`
	IsSystem        bool                            `gorm:"not null;default:false;index" json:"is_system"`
	Label           string                          `gorm:"type:varchar(250);not null" json:"label"`
	Placeholder     string                          `gorm:"type:varchar(250);not null" json:"placeholder"`
	AppliesTo       ProductAttributeAppliesToEnum   `gorm:"type:varchar(25);not null;check(applies_to IN ('product','variant'))" json:"applies_to"`
	DataType        ProductAttributeDataTypeEnum    `gorm:"type:varchar(25);not null;check(data_type IN ('string','int','float','bool','date','json','option','options'))" json:"data_type"`
	FieldType       ProductAttributeDisplayTypeEnum `gorm:"type:varchar(20);not null;default:'text'" json:"field_type"`
	ValidationRules datatypes.JSONSlice[string]     `gorm:"null" json:"validation_rules"`
	CreatedAt       time.Time                       `json:"created_at"`
	UpdatedAt       time.Time                       `json:"updated_at"`
}

func (e *ProductAttributeDefinition) TableName() string {
	return "product_attribute_definitions"
}

func (e *ProductAttributeDefinition) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
