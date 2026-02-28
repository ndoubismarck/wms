package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type TaskUserAssignment struct {
	ID        string    `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	TaskID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_task_user_assignments_task_user" json:"task_id"`
	UserID    string    `gorm:"type:varchar(150);not null;index;uniqueIndex:idx_task_user_assignments_task_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Task      *Task     `gorm:"foreignKey:TaskID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"task,omitempty"`
	User      *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func (e *TaskUserAssignment) TableName() string {
	return "task_user_assignments"
}

func (e *TaskUserAssignment) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
