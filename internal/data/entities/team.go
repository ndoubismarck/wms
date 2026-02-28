package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID  string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_teams_location_name" json:"location_id"`
	Name        string    `gorm:"type:varchar(150);not null;uniqueIndex:idx_teams_location_name" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	IsSystem    bool      `gorm:"not null;default:false;index" json:"is_system"`
	UserIDs     []string  `gorm:"-" json:"user_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Location    *Location `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
}

func (e *Team) TableName() string {
	return "teams"
}

func (e *Team) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
