package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductBrand struct {
	ID        string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationD string    `gorm:"type:varchar(150);index" json:"location_id"`
	Code      *string   `gorm:"type:varchar(10);uniqueIndex" json:"code,omitempty"`
	Name      string    `gorm:"type:varchar(150)" json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *ProductBrand) TableName() string {
	return "product_brands"
}

func (e *ProductBrand) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
