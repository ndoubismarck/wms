package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductAttributeRangeValue struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
type ProductAttribute struct {
	ID                           string                                          `gorm:"type:varchar(150);primaryKey" json:"id"`
	ProductID                    *string                                         `gorm:"type:varchar(150);index;null" json:"-"`
	ProductVariantID             *string                                         `gorm:"type:varchar(150);index;null" json:"-"`
	ProductAttributeDefinitionID string                                          `gorm:"type:varchar(150);not null;index" json:"-"`
	ProductAttributeOptionID     *string                                         `gorm:"type:varchar(150);null" json:"-"`
	Value                        any                                             `gorm:"-" json:"value"`
	ValueType                    ProductAttributeDataTypeEnum                    `gorm:"type:varchar(150);index;null" json:"data_type"`
	ValueString                  *string                                         `gorm:"type:text;null" json:"-"`
	ValueInt                     *int64                                          `gorm:"null" json:"-"`
	ValueFloat                   *float64                                        `gorm:"null" json:"-"`
	ValueBool                    *bool                                           `gorm:"null" json:"-"`
	ValueDate                    *time.Time                                      `gorm:"null" json:"-"`
	ValueJSON                    *datatypes.JSON                                 `gorm:"null" json:"-"`
	ValueRange                   *datatypes.JSONType[ProductAttributeRangeValue] `gorm:"null" json:"-"`
	ValueOption                  *ProductAttributeOption                         `gorm:"foreignKey:ProductAttributeOptionID;references:ID" json:"-"`
	ValueOptions                 []ProductAttributeOption                        `gorm:"foreignKey:ProductAttributeID;references:ID" json:"-"`
	CreatedAt                    time.Time                                       `json:"created_at"`
	UpdatedAt                    time.Time                                       `json:"updated_at"`
	Product                      *Product                                        `gorm:"foreignKey:ProductID;references:ID" json:"-"`
	ProductVariant               *ProductVariant                                 `gorm:"foreignKey:ProductVariantID;references:ID" json:"-"`
	ProductAttributeDefinition   *ProductAttributeDefinition                     `gorm:"foreignKey:ProductAttributeDefinitionID;references:ID" json:"-"`
}

func (e *ProductAttribute) TableName() string {
	return "product_attributes"
}

func (e *ProductAttribute) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}

func (e *ProductAttribute) AfterFind(tx *gorm.DB) error {
	if e.ValueString != nil {
		e.Value = *e.ValueString
		return nil
	}
	if e.ValueInt != nil {
		e.Value = *e.ValueInt
		return nil
	}
	if e.ValueFloat != nil {
		e.Value = *e.ValueFloat
		return nil
	}
	if e.ValueBool != nil {
		e.Value = *e.ValueBool
		return nil
	}
	if e.ValueDate != nil {
		e.Value = *e.ValueDate
		return nil
	}
	if e.ValueJSON != nil {
		e.Value = *e.ValueJSON
		return nil
	}
	if e.ValueRange != nil {
		e.Value = e.ValueRange.Data()
		return nil
	}
	if len(e.ValueOptions) > 0 {
		e.Value = e.ValueOptions
		return nil
	}
	if e.ValueOption != nil {
		e.Value = *e.ValueOption
		return nil
	}
	if e.ProductAttributeOptionID != nil {
		e.Value = *e.ProductAttributeOptionID
	}
	return nil
}
