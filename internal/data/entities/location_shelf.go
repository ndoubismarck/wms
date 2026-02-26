package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type LocationShelf struct {
	ID          string       `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	BayID       string       `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_location_shelves_bay_code" json:"bay_id"`
	Code        string       `gorm:"type:varchar(10);not null;uniqueIndex:idx_location_shelves_bay_code" json:"code"`
	Name        string       `gorm:"type:varchar(150)" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
	Location    string       `gorm:"->" json:"location"`
	Bay         *LocationBay `gorm:"foreignKey:BayID;references:ID" json:"bay,omitempty"`
	UpdatedAt   time.Time    `json:"updated_at"`
	CreatedAt   time.Time    `json:"created_at"`
}

func (e *LocationShelf) TableName() string {
	return "locations_shelves"
}

func (e *LocationShelf) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
