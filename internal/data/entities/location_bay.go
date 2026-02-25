package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type LocationBay struct {
	ID          string         `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	AisleID     string         `gorm:"type:varchar(150)" json:"aisle_id"`
	Code        string         `gorm:"type:varchar(10);uniqueIndex" json:"code"`
	Name        string         `gorm:"type:varchar(150)" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Location    string         `gorm:"->" json:"location"`
	Aisle       *LocationAisle `gorm:"foreignKey:AisleID;references:ID" json:"aisle,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

func (e *LocationBay) TableName() string {
	return "location_bays"
}

func (e *LocationBay) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
