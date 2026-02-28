package tasks

import (
	"server/internal/core/pagination"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database"
	"server/internal/services/providers/database/queries"
	"strings"
)

type Service struct {
	ctx       types.IContext
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	return &Service{
		ctx:       ctx,
		providers: providers,
	}
}

func normalizeIDs(values []string) []string {
	result := []string{}
	index := map[string]bool{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if index[trimmed] {
			continue
		}
		index[trimmed] = true
		result = append(result, trimmed)
	}
	return result
}

func validPriority(priority entities.TaskPriorityEnum) bool {
	switch priority {
	case entities.TaskPriorityLow, entities.TaskPriorityNormal, entities.TaskPriorityHigh, entities.TaskPriorityUrgent:
		return true
	default:
		return false
	}
}

func validStatus(status entities.TaskStatusEnum) bool {
	switch status {
	case entities.TaskStatusPending, entities.TaskStatusInProgress, entities.TaskStatusCompleted, entities.TaskStatusCancelled:
		return true
	default:
		return false
	}
}

func (s *Service) ensureScheduler(query *database.Query, scheduledByUserID *string, system bool) (entities.TaskScheduledByTypeEnum, *string, types.ServiceResultCode, error) {
	if system {
		return entities.TaskScheduledBySystem, nil, types.ServiceResultCodeSuccess, nil
	}
	if scheduledByUserID == nil || strings.TrimSpace(*scheduledByUserID) == "" {
		return "", nil, types.ServiceResultCodeForbidden, nil
	}
	userID := strings.TrimSpace(*scheduledByUserID)
	user, exists, err := query.Users().FindOneByID(userID)
	if err != nil {
		return "", nil, types.ServiceResultCodeInternal, err
	}
	if !exists || user.Role != entities.UserRoleAdmin {
		return "", nil, types.ServiceResultCodeForbidden, nil
	}
	return entities.TaskScheduledByUser, &userID, types.ServiceResultCodeSuccess, nil
}

func (s *Service) validateUserIDs(query *database.Query, userIDs []string) (types.ValidationResult, []string, error) {
	cleanIDs := normalizeIDs(userIDs)
	for _, userID := range cleanIDs {
		exists, err := query.Users().ExistsByID(userID)
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			return types.ValidationResult{"user_ids": {"contains unknown user IDs"}}, nil, nil
		}
	}
	return nil, cleanIDs, nil
}

func (s *Service) validateTeamIDs(query *database.Query, locationID string, teamIDs []string) (types.ValidationResult, []string, error) {
	cleanIDs := normalizeIDs(teamIDs)
	for _, teamID := range cleanIDs {
		exists, err := query.Teams().Exists(queries.TeamsParams{ID: &teamID, LocationID: &locationID})
		if err != nil {
			return nil, nil, err
		}
		if !exists {
			return types.ValidationResult{"team_ids": {"contains unknown or out-of-location team IDs"}}, nil, nil
		}
	}
	return nil, cleanIDs, nil
}

func (s *Service) setTaskAssignments(query *database.Query, taskID string, userIDs []string, teamIDs []string) error {
	if err := query.TaskUserAssignments().DeleteByTaskID(taskID); err != nil {
		return err
	}
	if err := query.TaskTeamAssignments().DeleteByTaskID(taskID); err != nil {
		return err
	}
	userAssignments := []entities.TaskUserAssignment{}
	for _, userID := range userIDs {
		userAssignments = append(userAssignments, entities.TaskUserAssignment{TaskID: taskID, UserID: userID})
	}
	if err := query.TaskUserAssignments().CreateMany(userAssignments); err != nil {
		return err
	}
	teamAssignments := []entities.TaskTeamAssignment{}
	for _, teamID := range teamIDs {
		teamAssignments = append(teamAssignments, entities.TaskTeamAssignment{TaskID: taskID, TeamID: teamID})
	}
	if err := query.TaskTeamAssignments().CreateMany(teamAssignments); err != nil {
		return err
	}
	return nil
}

func (s *Service) enrichTask(query *database.Query, task entities.Task) (entities.Task, error) {
	taskID := task.ID
	userIDs, err := query.TaskUserAssignments().FindUserIDs(queries.TaskUserAssignmentsParams{TaskID: &taskID})
	if err != nil {
		return entities.Task{}, err
	}
	teamIDs, err := query.TaskTeamAssignments().FindTeamIDs(queries.TaskTeamAssignmentsParams{TaskID: &taskID})
	if err != nil {
		return entities.Task{}, err
	}
	task.UserIDs = normalizeIDs(userIDs)
	task.TeamIDs = normalizeIDs(teamIDs)
	return task, nil
}

func (s *Service) Add(data AddData) (*AddResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &AddResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.Locations().Exists(queries.LocationsParams{ID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &AddResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	scheduledByType, scheduledByUserID, code, err := s.ensureScheduler(query, data.ScheduledByUserID, data.ScheduledBySystem)
	if err != nil {
		return nil, err
	}
	if code != types.ServiceResultCodeSuccess {
		return &AddResult{Code: code}, nil
	}
	priority := data.Priority
	if priority == "" {
		priority = entities.TaskPriorityNormal
	}
	if !validPriority(priority) {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"priority": {"invalid value"}}}, nil
	}
	status := data.Status
	if status == "" {
		status = entities.TaskStatusPending
	}
	if !validStatus(status) {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
	}
	userValidation, userIDs, err := s.validateUserIDs(query, data.UserIDs)
	if err != nil {
		return nil, err
	}
	if userValidation != nil {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: userValidation}, nil
	}
	teamValidation, teamIDs, err := s.validateTeamIDs(query, locationID, data.TeamIDs)
	if err != nil {
		return nil, err
	}
	if teamValidation != nil {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: teamValidation}, nil
	}
	title := strings.TrimSpace(data.Title)
	if title == "" {
		return &AddResult{
			Code: types.ServiceResultCodeInvalid,
			Validation: types.ValidationResult{
				"title": {"cannot be empty"},
			},
		}, nil
	}
	task, err := query.Tasks().Create(entities.Task{
		LocationID:        locationID,
		Title:             title,
		Description:       strings.TrimSpace(data.Description),
		Priority:          priority,
		Status:            status,
		DueAt:             data.DueAt,
		ScheduledByType:   scheduledByType,
		ScheduledByUserID: scheduledByUserID,
	})
	if err != nil {
		return nil, err
	}
	if err := s.setTaskAssignments(query, task.ID, userIDs, teamIDs); err != nil {
		return nil, err
	}
	task, err = s.enrichTask(query, task)
	if err != nil {
		return nil, err
	}
	return &AddResult{Code: types.ServiceResultCodeSuccess, Payload: AddResultPayload{Task: task}}, nil
}

func (s *Service) Update(data UpdateData) (*UpdateResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &UpdateResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	_, exists, err := query.Tasks().FindOne(queries.TasksParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	_, _, code, err := s.ensureScheduler(query, data.ScheduledByUserID, data.ScheduledBySystem)
	if err != nil {
		return nil, err
	}
	if code != types.ServiceResultCodeSuccess {
		return &UpdateResult{Code: code}, nil
	}
	hasChanges := false
	update := entities.Task{}
	if data.Title != nil {
		hasChanges = true
		title := strings.TrimSpace(*data.Title)
		if title == "" {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"title": {"cannot be empty"}}}, nil
		}
		update.Title = title
	}
	if data.Description != nil {
		hasChanges = true
		update.Description = strings.TrimSpace(*data.Description)
	}
	if data.Priority != nil {
		hasChanges = true
		if !validPriority(*data.Priority) {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"priority": {"invalid value"}}}, nil
		}
		update.Priority = *data.Priority
	}
	if data.Status != nil {
		hasChanges = true
		if !validStatus(*data.Status) {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"status": {"invalid value"}}}, nil
		}
		update.Status = *data.Status
	}
	if data.DueAt != nil {
		hasChanges = true
		update.DueAt = data.DueAt
	}
	if !hasChanges && data.UserIDs == nil && data.TeamIDs == nil {
		return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}
	if hasChanges {
		_, updated, err := query.Tasks().Update(queries.TasksParams{ID: &data.ID, LocationID: &locationID}, update)
		if err != nil {
			return nil, err
		}
		if !updated {
			return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
		}
	}
	taskID := data.ID
	userIDs, err := query.TaskUserAssignments().FindUserIDs(queries.TaskUserAssignmentsParams{TaskID: &taskID})
	if err != nil {
		return nil, err
	}
	teamIDs, err := query.TaskTeamAssignments().FindTeamIDs(queries.TaskTeamAssignmentsParams{TaskID: &taskID})
	if err != nil {
		return nil, err
	}
	if data.UserIDs != nil {
		userValidation, validated, err := s.validateUserIDs(query, data.UserIDs)
		if err != nil {
			return nil, err
		}
		if userValidation != nil {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: userValidation}, nil
		}
		userIDs = validated
	}
	if data.TeamIDs != nil {
		teamValidation, validated, err := s.validateTeamIDs(query, locationID, data.TeamIDs)
		if err != nil {
			return nil, err
		}
		if teamValidation != nil {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: teamValidation}, nil
		}
		teamIDs = validated
	}
	if data.UserIDs != nil || data.TeamIDs != nil {
		if err := s.setTaskAssignments(query, data.ID, normalizeIDs(userIDs), normalizeIDs(teamIDs)); err != nil {
			return nil, err
		}
	}
	task, exists, err := query.Tasks().FindOne(queries.TasksParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	task, err = s.enrichTask(query, task)
	if err != nil {
		return nil, err
	}
	return &UpdateResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateResultPayload{Task: task}}, nil
}

func (s *Service) Delete(data DeleteData) (*DeleteResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &DeleteResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &DeleteResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	exists, err := query.Tasks().Exists(queries.TasksParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	deleted, err := query.Tasks().Delete(queries.TasksParams{ID: &data.ID})
	if err != nil {
		return nil, err
	}
	if !deleted {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	return &DeleteResult{Code: types.ServiceResultCodeSuccess}, nil
}

func (s *Service) Get(data GetData) (*GetResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetResult{Code: types.ServiceResultCodeFailed}, nil
	}
	result, exists, err := query.Tasks().FindOne(queries.TasksParams{ID: data.ID, LocationID: &data.LocationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err = s.enrichTask(query, result)
	if err != nil {
		return nil, err
	}
	return &GetResult{Code: types.ServiceResultCodeSuccess, Payload: GetResultPayload{Task: result}}, nil
}

func (s *Service) GetMany(data GetManyData) (*GetManyResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetManyResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	results, paging, err := query.Tasks().FindMany(
		queries.TasksParams{LocationID: &locationID, Statuses: data.Statuses, Priorities: data.Priorities},
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}
	for idx := range results {
		task, err := s.enrichTask(query, results[idx])
		if err != nil {
			return nil, err
		}
		results[idx] = task
	}
	return &GetManyResult{Code: types.ServiceResultCodeSuccess, Payload: GetManyResultPayload{Tasks: results}, Pagination: paging}, nil
}
