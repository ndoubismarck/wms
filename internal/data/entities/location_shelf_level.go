package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type LocationShelfLevel struct {
	ID          string         `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	ShelfID     string         `gorm:"type:varchar(150)" json:"shelf_id"`
	Code        string         `gorm:"type:varchar(10);uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(150);unique" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"->" json:"location"`
	Shelf       *LocationShelf `gorm:"foreignKey:ShelfID;references:ID" json:"shelf,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (e *LocationShelfLevel) TableName() string {
	return "location_shelf_levels"
}

func (e *LocationShelfLevel) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
