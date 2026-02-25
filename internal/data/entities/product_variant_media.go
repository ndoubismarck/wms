package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductFileTypeEnum string

const (
	ProductFileTypeImage    ProductFileTypeEnum = "image"
	ProductFileTypeVideo    ProductFileTypeEnum = "video"
	ProductFileTypeDocument ProductFileTypeEnum = "document"
)

type ProductVariantMedia struct {
	ID               string              `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	ProductID        string              `gorm:"type:varchar(150)" json:"product_id"`
	ProductVariantID string              `gorm:"type:varchar(150)" json:"product_variant_id"`
	Path             string              `gorm:"type:varchar(1050);unique" json:"path"`
	FileType         ProductFileTypeEnum `gorm:"type:varchar(150);check(role IN ('image', 'video', 'document'))" json:"file_type"`
	FileMimeType     string              `gorm:"type:varchar(250)" json:"file_mime_type"`
	UpdatedAt        time.Time           `json:"updated_at"`
	CreatedAt        time.Time           `json:"created_at"`
}

func (e *ProductVariantMedia) TableName() string {
	return "product_variant_media"
}

func (e *ProductVariantMedia) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
