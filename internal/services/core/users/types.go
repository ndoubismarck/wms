package users

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
)

type (
	AddData struct {
		LocationID   string                  `json:"-" validate:"required"`
		Role         entities.UserRoleEnum   `json:"role"`
		Status       entities.UserStatusEnum `json:"status"`
		FirstName    string                  `json:"first_name"`
		LastName     string                  `json:"last_name"`
		PhoneNumber  string                  `json:"phone_number"`
		EmailAddress string                  `json:"email_address" validate:"required"`
		Password     string                  `json:"password" validate:"required"`
	}
	AddResult struct {
		Code       types.ServiceResultCode
		Payload    AddResultPayload
		Validation types.ValidationResult
	}
	AddResultPayload struct {
		User entities.User `json:"user"`
	}
)

type (
	UpdateData struct {
		LocationID   string                   `json:"-" validate:"required"`
		ID           string                   `json:"id" validate:"required"`
		Role         *entities.UserRoleEnum   `json:"role,omitempty"`
		Status       *entities.UserStatusEnum `json:"status,omitempty"`
		FirstName    *string                  `json:"first_name,omitempty"`
		LastName     *string                  `json:"last_name,omitempty"`
		PhoneNumber  *string                  `json:"phone_number,omitempty"`
		EmailAddress *string                  `json:"email_address,omitempty"`
		Password     *string                  `json:"password,omitempty"`
	}
	UpdateResult struct {
		Code       types.ServiceResultCode
		Payload    UpdateResultPayload
		Validation types.ValidationResult
	}
	UpdateResultPayload struct {
		User entities.User `json:"user"`
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
	GetManyData struct {
		types.PaginationParams
		LocationID string   `json:"-" form:"-" validate:"required"`
		Roles      []string `json:"roles" form:"roles"`
	}
	GetManyResult struct {
		Code       types.ServiceResultCode
		Payload    GetManyResultPayload
		Pagination *types.PaginationResult
	}
	GetManyResultPayload struct {
		Users []entities.User `json:"users"`
	}
)
