package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductAttributeOption struct {
	ID                 string            `gorm:"type:varchar(150);primaryKey" json:"id"`
	ProductAttributeID string            `gorm:"type:varchar(150);not null;index" json:"product_attribute_id"`
	UpdatedAt          time.Time         `json:"updated_at"`
	CreatedAt          time.Time         `json:"created_at"`
	ProductAttribute   *ProductAttribute `gorm:"foreignKey:ProductAttributeID;references:ID" json:"product_attribute,omitempty"`
}

func (e *ProductAttributeOption) TableName() string {
	return "product_attribute_options"
}

func (e *ProductAttributeOption) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
