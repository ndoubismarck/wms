package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type LocationBin struct {
	ID           string              `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	ShelfLevelID string              `gorm:"type:varchar(150)" json:"shelf_level_id"`
	Code         string              `gorm:"type:varchar(10);uniqueIndex" json:"code"`
	Name         string              `gorm:"type:varchar(150)" json:"name"`
	Description  string              `gorm:"type:text" json:"description"`
	Location     string              `gorm:"->" json:"location"`
	StockCount   uint64              `gorm:"-" json:"stock_count"`
	ShelfLevel   *LocationShelfLevel `gorm:"foreignKey:ShelfLevelID;references:ID" json:"shelf_level,omitempty"`
	UpdatedAt    time.Time           `json:"updated_at"`
	CreatedAt    time.Time           `json:"created_at"`
}

func (e *LocationBin) TableName() string {
	return "location_bins"
}

func (e *LocationBin) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
