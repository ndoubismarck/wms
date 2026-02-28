package entities

import (
	"server/internal/core/shared/utils/cryptoutil"
	"server/internal/core/shared/utils/stringutil"
	"time"

	"gorm.io/gorm"
)

type TaskPriorityEnum string

type TaskStatusEnum string

type TaskScheduledByTypeEnum string

const (
	TaskPriorityLow    TaskPriorityEnum = "low"
	TaskPriorityNormal TaskPriorityEnum = "normal"
	TaskPriorityHigh   TaskPriorityEnum = "high"
	TaskPriorityUrgent TaskPriorityEnum = "urgent"
)

const (
	TaskStatusPending    TaskStatusEnum = "pending"
	TaskStatusInProgress TaskStatusEnum = "in_progress"
	TaskStatusCompleted  TaskStatusEnum = "completed"
	TaskStatusCancelled  TaskStatusEnum = "cancelled"
)

const (
	TaskScheduledByUser   TaskScheduledByTypeEnum = "user"
	TaskScheduledBySystem TaskScheduledByTypeEnum = "system"
)

type Task struct {
	ID                string                  `gorm:"type:varchar(150);primaryKey;unique" json:"id"`
	LocationID        string                  `gorm:"type:varchar(150);not null;index" json:"location_id"`
	Title             string                  `gorm:"type:varchar(200);not null" json:"title"`
	Description       string                  `gorm:"type:text" json:"description"`
	Priority          TaskPriorityEnum        `gorm:"type:varchar(25);not null;default:'normal';check(priority IN ('low','normal','high','urgent'))" json:"priority"`
	Status            TaskStatusEnum          `gorm:"type:varchar(25);not null;default:'pending';check(status IN ('pending','in_progress','completed','cancelled'))" json:"status"`
	ScheduledByType   TaskScheduledByTypeEnum `gorm:"type:varchar(25);not null;default:'user';check(scheduled_by_type IN ('user','system'))" json:"scheduled_by_type"`
	ScheduledByUserID *string                 `gorm:"type:varchar(150);index;null" json:"scheduled_by_user_id,omitempty"`
	DueAt             *time.Time              `gorm:"null" json:"due_at,omitempty"`
	CreatedAt         time.Time               `json:"created_at"`
	UpdatedAt         time.Time               `json:"updated_at"`
	Location          *Location               `gorm:"foreignKey:LocationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"location,omitempty"`
	ScheduledByUser   *User                   `gorm:"foreignKey:ScheduledByUserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"scheduled_by_user,omitempty"`
	UserIDs           []string                `gorm:"-" json:"user_ids"`
	TeamIDs           []string                `gorm:"-" json:"team_ids"`
}

func (e *Task) TableName() string {
	return "tasks"
}

func (e *Task) BeforeCreate(tx *gorm.DB) error {
	if stringutil.IsEmpty(e.ID) {
		e.ID = cryptoutil.UID()
	}
	return nil
}
