package entities

import (
	"errors"
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID          string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	Code        uint32    `gorm:"uniqueIndex" json:"-"`
	Name        string    `gorm:"type:varchar(150);unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e *Location) TableName() string {
	return "locations"
}

func (e *Location) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	if e.Code == 0 {
		var data Location
		query := tx.Model(&Location{})
		if err := query.Order("code DESC").First(&data).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
		e.Code = data.Code + 1
	}
	return nil
}
