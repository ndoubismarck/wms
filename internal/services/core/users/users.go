package users

import (
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/securityutil"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database/queries"
	"strings"
)

type Service struct {
	ctx types.IContext
	tmp struct {
		admin *entities.User
	}
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func (s *Service) Add(data AddData) (*AddResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddResult{Code: types.ServiceResultCodeFailed}, nil
	}
	email := strings.TrimSpace(data.EmailAddress)
	exists, err := query.Users().ExistsByEmailAddress(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return &AddResult{
			Code: types.ServiceResultCodeInvalid,
			Validation: types.ValidationResult{
				"email_address": {"already exists"},
			},
		}, nil
	}
	password, err := securityutil.HashPassword(strings.TrimSpace(data.Password))
	if err != nil {
		return nil, err
	}
	role := data.Role
	if role == "" {
		role = entities.UserRoleSubscriber
	}
	status := data.Status
	if status == "" {
		status = entities.UserStatusActive
	}
	user, err := query.Users().Create(entities.User{
		Role:         role,
		Status:       status,
		FirstName:    strings.TrimSpace(data.FirstName),
		LastName:     strings.TrimSpace(data.LastName),
		PhoneNumber:  strings.TrimSpace(data.PhoneNumber),
		EmailAddress: email,
		Password:     password,
	})
	if err != nil {
		return nil, err
	}
	return &AddResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: AddResultPayload{
			User: user,
		},
	}, nil
}

func (s *Service) Update(data UpdateData) (*UpdateResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	hasChanges := false
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateResult{Code: types.ServiceResultCodeFailed}, nil
	}
	update := entities.User{}
	if data.Role != nil {
		hasChanges = true
		update.Role = *data.Role
	}
	if data.Status != nil {
		hasChanges = true
		update.Status = *data.Status
	}
	if data.FirstName != nil {
		hasChanges = true
		update.FirstName = strings.TrimSpace(*data.FirstName)
	}
	if data.LastName != nil {
		hasChanges = true
		update.LastName = strings.TrimSpace(*data.LastName)
	}
	if data.PhoneNumber != nil {
		hasChanges = true
		update.PhoneNumber = strings.TrimSpace(*data.PhoneNumber)
	}
	if data.EmailAddress != nil {
		hasChanges = true
		email := strings.TrimSpace(*data.EmailAddress)
		if email == "" {
			return &UpdateResult{
				Code: types.ServiceResultCodeInvalid,
				Validation: types.ValidationResult{
					"email_address": {"cannot be empty"},
				},
			}, nil
		}
		existing, found, err := query.Users().FindOneByEmailAddress(email)
		if err != nil {
			return nil, err
		}
		if found && existing.ID != data.ID {
			return &UpdateResult{
				Code: types.ServiceResultCodeInvalid,
				Validation: types.ValidationResult{
					"email_address": {"already exists"},
				},
			}, nil
		}
		update.EmailAddress = email
	}
	if data.Password != nil {
		hasChanges = true
		password := strings.TrimSpace(*data.Password)
		if password == "" {
			return &UpdateResult{
				Code: types.ServiceResultCodeInvalid,
				Validation: types.ValidationResult{
					"password": {"cannot be empty"},
				},
			}, nil
		}
		hash, err := securityutil.HashPassword(password)
		if err != nil {
			return nil, err
		}
		update.Password = hash
	}
	if !hasChanges {
		return &UpdateResult{
			Code: types.ServiceResultCodeInvalid,
			Validation: types.ValidationResult{
				"payload": {"no fields to update"},
			},
		}, nil
	}
	_, exists, err := query.Users().UpdateByID(data.ID, update)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	user, exists, err := query.Users().FindOneByID(data.ID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &UpdateResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: UpdateResultPayload{
			User: user,
		},
	}, nil
}

func (s *Service) Delete(data DeleteData) (*DeleteResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteResult{Code: types.ServiceResultCodeFailed}, nil
	}
	deleted, err := query.Users().DeleteByID(data.ID)
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) GetMany(data GetManyData) (*GetManyResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetManyResult{Code: types.ServiceResultCodeFailed}, nil
	}
	results, paging, err := query.Users().FindAll(
		queries.UserFilters{Roles: data.Roles},
		pagination.New(types.PaginationParams{
			Page:  data.Page,
			Limit: data.Limit,
			Order: data.Order,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &GetManyResult{
		Code: types.ServiceResultCodeSuccess,
		Payload: GetManyResultPayload{
			Users: results,
		},
		Pagination: paging,
	}, nil
}
