package entities

import (
	"fmt"
	"math/rand"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type ProductVariant struct {
	ID           string                `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	ProductID    string                `gorm:"type:varchar(150);index" json:"product_id"`
	SKU          string                `gorm:"type:varchar(50);uniqueIndex" json:"sku"`
	SellingPrice uint64                `json:"selling_price"`
	UpdatedAt    time.Time             `json:"updated_at"`
	CreatedAt    time.Time             `json:"created_at"`
	Product      *Product              `gorm:"foreignKey:ProductID;references:ID" json:"product,omitempty"`
	Attributes   []ProductAttribute    `gorm:"foreignKey:ProductVariantID;references:ID" json:"attributes,omitempty"`
	Media        []ProductVariantMedia `gorm:"foreignKey:ProductVariantID;references:ID" json:"media,omitempty"`
}

func (e *ProductVariant) TableName() string {
	return "product_variants"
}

func (e *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	if len(e.SKU) == 0 {
		sku, err := e.generateSku(tx)
		if err != nil {
			return err
		}
		e.SKU = sku
	}
	return nil
}

func (e *ProductVariant) generateSku(tx *gorm.DB) (string, error) {
	const maxRetries = 50
	for i := range maxRetries {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i)))
		str := fmt.Sprintf("%09d", rnd.Intn(100_000_000))
		sku := fmt.Sprintf("%d", stringutil.ToUint64(str))
		var count int64
		if err := tx.Model(&ProductVariant{}).Where("sku = ?", sku).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return sku, nil
		}
	}
	return "", fmt.Errorf("failed to generate unique SKU after %d attempts", maxRetries)
}
