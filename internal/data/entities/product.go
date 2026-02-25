package entities

import (
	"html"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductStatusEnum string

const (
	ProductStatusDraft       ProductStatusEnum = "draft"
	ProductStatusPending     ProductStatusEnum = "pending"
	ProductStatusAvailable   ProductStatusEnum = "available"
	ProductStatusUnavailable ProductStatusEnum = "unavailable"
)

type Product struct {
	ID            string              `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationD     string              `gorm:"type:varchar(150);index" json:"location_id"`
	SubcategoryID string              `gorm:"type:varchar(150);index" json:"subcategory_id"`
	BrandID       string              `gorm:"type:varchar(255);index" json:"brand_id"`
	Name          string              `gorm:"type:varchar(250);index" json:"name"`
	Description   string              `gorm:"type:text;index" json:"description"`
	Status        ProductStatusEnum   `gorm:"type:varchar(150);index;check(status IN ('draft', 'pending', 'available', 'unavailable'))" json:"status"`
	UpdatedAt     time.Time           `json:"updated_at"`
	CreatedAt     time.Time           `json:"created_at"`
	Brand         *ProductBrand       `gorm:"foreignKey:BrandID;references:ID" json:"brand,omitempty"`
	Subcategory   *ProductSubcategory `gorm:"foreignKey:SubcategoryID;references:ID" json:"subcategory,omitempty"`
	Attributes    []ProductAttribute  `gorm:"foreignKey:ProductID;references:ID" json:"attributes,omitempty"`
	Variants      []ProductVariant    `gorm:"foreignKey:ProductID;references:ID" json:"variants,omitempty"`
}

func (e *Product) TableName() string {
	return "products"
}

func (e *Product) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return e.escapeValues()
}

func (e *Product) BeforeUpdate(tx *gorm.DB) error {
	return e.escapeValues()
}

func (e *Product) AfterFind(tx *gorm.DB) error {
	return e.unescapeValues()
}

func (e *Product) escapeValues() error {
	if stringutil.IsHTML(e.Description) {
		e.Description = html.EscapeString(e.Description)
	}
	return nil
}

func (e *Product) unescapeValues() error {
	e.Description = html.UnescapeString(e.Description)
	return nil
}
