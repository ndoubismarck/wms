package teams

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

func (s *Service) setTeamMembers(query *database.Query, teamID string, userIDs []string) error {
	if err := query.TeamMembers().DeleteByTeamID(teamID); err != nil {
		return err
	}
	assignments := []entities.TeamMember{}
	for _, userID := range userIDs {
		assignments = append(assignments, entities.TeamMember{TeamID: teamID, UserID: userID})
	}
	if err := query.TeamMembers().CreateMany(assignments); err != nil {
		return err
	}
	return nil
}

func (s *Service) enrichTeam(query *database.Query, team entities.Team) (entities.Team, error) {
	teamID := team.ID
	userIDs, err := query.TeamMembers().FindUserIDs(queries.TeamMembersParams{TeamID: &teamID})
	if err != nil {
		return entities.Team{}, err
	}
	team.UserIDs = normalizeIDs(userIDs)
	return team, nil
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
	name := strings.TrimSpace(data.Name)
	if name == "" {
		return &AddResult{
			Code: types.ServiceResultCodeInvalid,
			Validation: types.ValidationResult{
				"name": {"cannot be empty"},
			},
		}, nil
	}
	exists, err = query.Teams().Exists(queries.TeamsParams{LocationID: &locationID, Name: &name})
	if err != nil {
		return nil, err
	}
	if exists {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"already exists in this location"}}}, nil
	}
	userValidation, userIDs, err := s.validateUserIDs(query, data.UserIDs)
	if err != nil {
		return nil, err
	}
	if userValidation != nil {
		return &AddResult{Code: types.ServiceResultCodeInvalid, Validation: userValidation}, nil
	}
	team, err := query.Teams().Create(entities.Team{
		LocationID:  locationID,
		Name:        name,
		Description: strings.TrimSpace(data.Description),
	})
	if err != nil {
		return nil, err
	}
	if err := s.setTeamMembers(query, team.ID, userIDs); err != nil {
		return nil, err
	}
	team, err = s.enrichTeam(query, team)
	if err != nil {
		return nil, err
	}
	return &AddResult{Code: types.ServiceResultCodeSuccess, Payload: AddResultPayload{Team: team}}, nil
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
	current, exists, err := query.Teams().FindOne(queries.TeamsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	if current.IsSystem {
		return &UpdateResult{Code: types.ServiceResultCodeForbidden}, nil
	}
	hasChanges := false
	update := entities.Team{}
	if data.Name != nil {
		hasChanges = true
		name := strings.TrimSpace(*data.Name)
		if name == "" {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"cannot be empty"}}}, nil
		}
		existing, found, err := query.Teams().FindOne(queries.TeamsParams{LocationID: &locationID, Name: &name})
		if err != nil {
			return nil, err
		}
		if found && existing.ID != data.ID {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"name": {"already exists in this location"}}}, nil
		}
		update.Name = name
	}
	if data.Description != nil {
		hasChanges = true
		update.Description = strings.TrimSpace(*data.Description)
	}
	if !hasChanges && data.UserIDs == nil {
		return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: types.ValidationResult{"payload": {"no fields to update"}}}, nil
	}
	if hasChanges {
		_, updated, err := query.Teams().Update(queries.TeamsParams{ID: &data.ID, LocationID: &locationID}, update)
		if err != nil {
			return nil, err
		}
		if !updated {
			return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
		}
	}
	if data.UserIDs != nil {
		userValidation, userIDs, err := s.validateUserIDs(query, data.UserIDs)
		if err != nil {
			return nil, err
		}
		if userValidation != nil {
			return &UpdateResult{Code: types.ServiceResultCodeInvalid, Validation: userValidation}, nil
		}
		if err := s.setTeamMembers(query, data.ID, userIDs); err != nil {
			return nil, err
		}
	}
	team, exists, err := query.Teams().FindOne(queries.TeamsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &UpdateResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	team, err = s.enrichTeam(query, team)
	if err != nil {
		return nil, err
	}
	return &UpdateResult{Code: types.ServiceResultCodeSuccess, Payload: UpdateResultPayload{Team: team}}, nil
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
	team, exists, err := query.Teams().FindOne(queries.TeamsParams{ID: &data.ID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &DeleteResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	if team.IsSystem {
		return &DeleteResult{Code: types.ServiceResultCodeForbidden}, nil
	}
	deleted, err := query.Teams().Delete(queries.TeamsParams{ID: &data.ID})
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
	result, exists, err := query.Teams().FindOne(queries.TeamsParams{ID: data.ID, LocationID: &data.LocationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	result, err = s.enrichTeam(query, result)
	if err != nil {
		return nil, err
	}
	return &GetResult{Code: types.ServiceResultCodeSuccess, Payload: GetResultPayload{Team: result}}, nil
}

func (s *Service) GetMany(data GetManyData) (*GetManyResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetManyResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	results, paging, err := query.Teams().FindMany(
		queries.TeamsParams{LocationID: &locationID},
		pagination.New(types.PaginationParams{Page: data.Page, Limit: data.Limit, Order: data.Order}),
	)
	if err != nil {
		return nil, err
	}
	for idx := range results {
		team, err := s.enrichTeam(query, results[idx])
		if err != nil {
			return nil, err
		}
		results[idx] = team
	}
	return &GetManyResult{Code: types.ServiceResultCodeSuccess, Payload: GetManyResultPayload{Teams: results}, Pagination: paging}, nil
}

func (s *Service) SetMembers(data SetMembersData) (*SetMembersResult, error) {
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &SetMembersResult{Code: types.ServiceResultCodeInvalid, Validation: vld}, nil
	}
	query, ok := s.providers.Database().Query()
	if !ok {
		return &SetMembersResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	teamID := data.TeamID
	exists, err := query.Teams().Exists(queries.TeamsParams{ID: &teamID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &SetMembersResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	validation, userIDs, err := s.validateUserIDs(query, data.UserIDs)
	if err != nil {
		return nil, err
	}
	if validation != nil {
		return &SetMembersResult{Code: types.ServiceResultCodeInvalid, Validation: validation}, nil
	}
	if err := s.setTeamMembers(query, data.TeamID, userIDs); err != nil {
		return nil, err
	}
	return &SetMembersResult{Code: types.ServiceResultCodeSuccess, Payload: SetMembersResultPayload{UserIDs: userIDs}}, nil
}

func (s *Service) GetMembers(data GetMembersData) (*GetMembersResult, error) {
	query, ok := s.providers.Database().Query()
	if !ok {
		return &GetMembersResult{Code: types.ServiceResultCodeFailed}, nil
	}
	locationID := data.LocationID
	teamID := data.TeamID
	exists, err := query.Teams().Exists(queries.TeamsParams{ID: &teamID, LocationID: &locationID})
	if err != nil {
		return nil, err
	}
	if !exists {
		return &GetMembersResult{Code: types.ServiceResultCodeNotFound}, nil
	}
	userIDs, err := query.TeamMembers().FindUserIDs(queries.TeamMembersParams{TeamID: &teamID})
	if err != nil {
		return nil, err
	}
	return &GetMembersResult{Code: types.ServiceResultCodeSuccess, Payload: GetMembersResultPayload{UserIDs: normalizeIDs(userIDs)}}, nil
}
