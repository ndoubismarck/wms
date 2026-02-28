package teams

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
)

type (
	AddData struct {
		LocationID  string   `json:"-" validate:"required"`
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		UserIDs     []string `json:"user_ids"`
	}
	AddResult struct {
		Code       types.ServiceResultCode
		Payload    AddResultPayload
		Validation types.ValidationResult
	}
	AddResultPayload struct {
		Team entities.Team `json:"team"`
	}
)

type (
	UpdateData struct {
		LocationID  string   `json:"-" validate:"required"`
		ID          string   `json:"id" validate:"required"`
		Name        *string  `json:"name,omitempty"`
		Description *string  `json:"description,omitempty"`
		UserIDs     []string `json:"user_ids,omitempty"`
	}
	UpdateResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateResultPayload
		Validation types.ValidationResult
	}
	UpdateResultPayload struct {
		Team entities.Team `json:"team"`
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
		Team entities.Team `json:"team"`
	}
)

type (
	GetManyData struct {
		types.PaginationParams
		LocationID string `json:"-" form:"-" validate:"required"`
	}
	GetManyResult struct {
		Code       types.ServiceResultCode
		Payload    GetManyResultPayload
		Pagination *types.PaginationResult
	}
	GetManyResultPayload struct {
		Teams []entities.Team `json:"teams"`
	}
)

type (
	SetMembersData struct {
		LocationID string   `json:"-" validate:"required"`
		TeamID     string   `json:"team_id" validate:"required"`
		UserIDs    []string `json:"user_ids"`
	}
	SetMembersResult struct {
		Code       types.ServiceResultCode
		Payload    SetMembersResultPayload
		Validation types.ValidationResult
	}
	SetMembersResultPayload struct {
		UserIDs []string `json:"user_ids"`
	}
)

type (
	GetMembersData struct {
		LocationID string `json:"-" form:"-" validate:"required"`
		TeamID     string `json:"team_id" form:"team_id" validate:"required"`
	}
	GetMembersResult struct {
		Code    types.ServiceResultCode
		Payload GetMembersResultPayload
	}
	GetMembersResultPayload struct {
		UserIDs []string `json:"user_ids"`
	}
)
