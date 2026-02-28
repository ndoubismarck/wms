package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type TaskTeamAssignment struct {
	ID        string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	TaskID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_task_team_assignments_task_team" json:"task_id"`
	TeamID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_task_team_assignments_task_team" json:"team_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Task      *Task     `gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"task,omitempty"`
	Team      *Team     `gorm:"foreignKey:TeamID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team,omitempty"`
}

func (e *TaskTeamAssignment) TableName() string {
	return "task_team_assignments"
}

func (e *TaskTeamAssignment) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
