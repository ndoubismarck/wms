package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type TeamMember struct {
	ID        string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	TeamID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_team_members_team_user" json:"team_id"`
	UserID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_team_members_team_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Team      *Team     `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team,omitempty"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func (e *TeamMember) TableName() string {
	return "team_members"
}

func (e *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
