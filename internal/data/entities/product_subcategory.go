package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductSubcategory struct {
	ID         string           `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	CategoryID string           `gorm:"type:varchar(150)" json:"category_id"`
	Code       *string          `gorm:"type:varchar(10);uniqueIndex" json:"code,omitempty"`
	Name       string           `gorm:"type:varchar(150)" json:"name"`
	UpdatedAt  time.Time        `json:"updated_at"`
	CreatedAt  time.Time        `json:"created_at"`
	Category   *ProductCategory `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
}

func (e *ProductSubcategory) TableName() string {
	return "product_subcategories"
}

func (e *ProductSubcategory) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
