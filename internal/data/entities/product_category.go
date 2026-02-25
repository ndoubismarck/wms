package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID        string        `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	BrandID   string        `gorm:"type:varchar(150);index" json:"brand_id"`
	Code      *string       `gorm:"type:varchar(10);uniqueIndex" json:"code,omitempty"`
	Name      string        `gorm:"type:varchar(150)" json:"name"`
	UpdatedAt time.Time     `json:"updated_at"`
	CreatedAt time.Time     `json:"created_at"`
	Brand     *ProductBrand `gorm:"foreignKey:BrandID;references:ID" json:"brand,omitempty"`
}

func (e *ProductCategory) TableName() string {
	return "product_categories"
}

func (e *ProductCategory) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
