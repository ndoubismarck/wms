package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type LocationAisle struct {
	ID          string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID  string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_location_aisles_location_code" json:"location_id"`
	Code        string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_location_aisles_location_code" json:"code"`
	Name        string    `gorm:"type:varchar(150);unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Location    *Location `gorm:"foreignKey:LocationID;references:ID" json:"location,omitempty"`
}

func (e *LocationAisle) TableName() string {
	return "location_aisles"
}

func (e *LocationAisle) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
