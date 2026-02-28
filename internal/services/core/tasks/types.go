package tasks

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"time"
)

type (
	AddData struct {
		LocationID        string                    `json:"-" validate:"required"`
		Title             string                    `json:"title" validate:"required"`
		Description       string                    `json:"description"`
		Priority          entities.TaskPriorityEnum `json:"priority"`
		Status            entities.TaskStatusEnum   `json:"status"`
		DueAt             *time.Time                `json:"due_at,omitempty"`
		UserIDs           []string                  `json:"user_ids"`
		TeamIDs           []string                  `json:"team_ids"`
		ScheduledByUserID *string                   `json:"-"`
		ScheduledBySystem bool                      `json:"-"`
	}
	AddResult struct {
		Code       types.ServiceResultCode
		Payload    AddResultPayload
		Validation types.ValidationResult
	}
	AddResultPayload struct {
		Task entities.Task `json:"task"`
	}
)

type (
	UpdateData struct {
		LocationID        string                     `json:"-" validate:"required"`
		ID                string                     `json:"id" validate:"required"`
		Title             *string                    `json:"title,omitempty"`
		Description       *string                    `json:"description,omitempty"`
		Priority          *entities.TaskPriorityEnum `json:"priority,omitempty"`
		Status            *entities.TaskStatusEnum   `json:"status,omitempty"`
		DueAt             *time.Time                 `json:"due_at,omitempty"`
		UserIDs           []string                   `json:"user_ids,omitempty"`
		TeamIDs           []string                   `json:"team_ids,omitempty"`
		ScheduledByUserID *string                    `json:"-"`
		ScheduledBySystem bool                       `json:"-"`
	}
	UpdateResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateResultPayload
		Validation types.ValidationResult
	}
	UpdateResultPayload struct {
		Task entities.Task `json:"task"`
	}
)

type (
	DeleteData struct {
		LocationID string `json:"-" validate:"required"`
		ID         string `json:"id" validate:"required"`
	}
	DeleteResult struct {
		Code       types.ServiceResultCode
		Validation types.ValidationResult
	}
)

type (
	GetData struct {
		LocationID string  `json:"-" form:"-" validate:"required"`
		ID         *string `json:"id" form:"id"`
	}
	GetResult struct {
		Code    types.ServiceResultCode
		Payload GetResultPayload
	}
	GetResultPayload struct {
		Task entities.Task `json:"task"`
	}
)

type (
	GetManyData struct {
		types.PaginationParams
		LocationID string                      `json:"-" form:"-" validate:"required"`
		Statuses   []entities.TaskStatusEnum   `json:"statuses" form:"statuses"`
		Priorities []entities.TaskPriorityEnum `json:"priorities" form:"priorities"`
	}
	GetManyResult struct {
		Code       types.ServiceResultCode
		Payload    GetManyResultPayload
		Pagination *types.PaginationResult
	}
	GetManyResultPayload struct {
		Tasks []entities.Task `json:"tasks"`
	}
)
