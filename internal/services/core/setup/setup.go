package setup

import (
	"fmt"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/securityutil"
	"server/internal/core/shared/utils/stringutil"
	"server/internal/data/entities"
	"server/internal/services/providers"
	"server/internal/services/providers/database"
	"server/internal/services/providers/database/queries"
	"strings"
)

type Service struct {
	ctx types.IContext
	tmp struct {
		admin    []entities.User
		location []entities.Location
	}
	providers *providers.Providers
}

func New(ctx types.IContext, providers *providers.Providers) *Service {
	s := &Service{
		ctx:       ctx,
		providers: providers,
	}
	ctx.Hooks().OnStarted(s.loadTmp)
	return s
}

func (s *Service) loadTmp() error {
	defer func(s *Service) {
		s.tmp.admin = nil
		s.tmp.location = nil
	}(s)
	query, ok := s.providers.Database().Query()
	if !ok {
		return nil
	}
	if s.tmp.admin != nil {
		for _, val := range s.tmp.admin {
			if err := s.createAdmin(query, val); err != nil {
				return err
			}
		}
		s.tmp.admin = []entities.User{}
	}
	if s.tmp.location != nil {
		for _, val := range s.tmp.location {
			if err := s.createLocation(query, &val); err != nil {
				return err
			}
		}
		s.tmp.location = []entities.Location{}
	}
	return nil
}

func (s *Service) adminSetup(query *database.Query) (bool, error) {
	count, err := query.Users().Count(queries.UserFilters{
		Roles: []string{entities.UserRoleAdmin},
	})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *Service) Get(data GetData) (*GetResult, error) {
	query, ok := s.providers.Database().Query()
	payload := GetResultPayload{}
	payload.Setup.Database = ok
	if ok {
		admin, err := s.adminSetup(query)
		if err != nil {
			return nil, err
		}
		payload.Setup.Admin = admin
		if data.User != nil && data.User.Role == entities.UserRoleAdmin {
			// todo: this should be added to or fetched from users's setup
			location, exists, err := query.Locations().FindOne(queries.LocationsParams{})
			if err != nil {
				return nil, err
			}
			if exists {
				payload.Setup.Location = &location
			}
		}
	}
	return &GetResult{
		Code:    types.ServiceResultCodeSuccess,
		Payload: payload,
	}, nil
}

func (s *Service) Create(data CreateData) (*CreateResult, error) {
	var (
		write  bool
		result CreateResult
	)
	vld, ok, err := s.providers.Validation().ValidateStruct(data)
	if err != nil {
		return nil, err
	}
	if !ok {
		return &CreateResult{
			Code:       types.ServiceResultCodeInvalid,
			Validation: vld,
		}, nil
	}
	cnf, err := s.ctx.Config().Get()
	if err != nil {
		return nil, err
	}
	query, ok := s.providers.Database().Query()
	if data.Location != nil {
		location := entities.Location{
			Name:        data.Location.Name,
			Description: data.Location.Name,
		}
		if ok {
			if err = s.createLocation(query, &location); err != nil {
				return nil, err
			}
		} else {
			if s.tmp.location == nil {
				s.tmp.location = []entities.Location{}
			}
			s.tmp.location = append(s.tmp.location, location)
		}
		result.Payload.Setup.Location = true
	}
	if data.Admin != nil {
		password, err := securityutil.HashPassword(data.Admin.Password)
		if err != nil {
			return nil, err
		}
		admin := entities.User{
			Role:         entities.UserRoleAdmin,
			Status:       entities.UserStatusActive,
			FirstName:    strings.TrimSpace(data.Admin.FirstName),
			LastName:     strings.TrimSpace(data.Admin.LastName),
			PhoneNumber:  strings.TrimSpace(data.Admin.PhoneNumber),
			EmailAddress: strings.TrimSpace(data.Admin.EmailAddress),
			Password:     password,
		}
		if ok {
			if err = s.createAdmin(query, admin); err != nil {
				return nil, err
			}
		} else {
			if s.tmp.admin == nil {
				s.tmp.admin = []entities.User{}
			}
			s.tmp.admin = append(s.tmp.admin, admin)
		}
		result.Payload.Setup.Admin = true
	}

	if data.Database != nil {
		if data.Database.SQLite != nil {
			cnf.Database.SQLite = &types.SqliteDatabaseConfig{
				EncryptionEnabled: data.Database.SQLite.EncryptionEnabled,
			}
			result.Payload.Setup.Database = true
			write = true
		}

		if data.Database.MariaDB != nil {
			cnf.Database.MariaDB = &types.MariaDBDatabaseConfig{
				Address:      fmt.Sprintf("%s:%d", data.Database.MariaDB.Host, data.Database.MariaDB.Port),
				Username:     data.Database.MariaDB.Username,
				Password:     data.Database.MariaDB.Password,
				DatabaseName: data.Database.MariaDB.DatabaseName,
			}
			result.Payload.Setup.Database = true
			write = true
		}
		if data.Database.Postgres != nil {
			cnf.Database.Postgres = &types.PostgresDatabaseConfig{
				Address:      fmt.Sprintf("%s:%d", data.Database.Postgres.Host, data.Database.Postgres.Port),
				Username:     data.Database.Postgres.Username,
				Password:     data.Database.Postgres.Password,
				DatabaseName: data.Database.Postgres.DatabaseName,
			}
			result.Payload.Setup.Database = true
			write = true
		}
	}
	if write {
		if err = s.ctx.Config().Set(cnf); err != nil {
			return nil, err
		}
	}
	result.Code = types.ServiceResultCodeSuccess
	return &result, nil
}

func (s *Service) createAdmin(query *database.Query, user entities.User) error {
	s.ctx.Logger().Debugf("upserting admin user %s...", user.EmailAddress)
	if _, err := query.Users().Upsert(user); err != nil {
		s.ctx.Logger().Errorf("could not upsert admin user %s: %v", user.EmailAddress, err)
		return err
	}
	s.ctx.Logger().Debugf("upserting admin user %s successful", user.EmailAddress)
	return nil
}

func (s *Service) createLocation(query *database.Query, location *entities.Location) error {
	currentLocation, ok, err := query.Locations().FindOne(queries.LocationsParams{
		ID: stringutil.ToPtr(location.ID),
	})
	if err != nil {
		s.ctx.Logger().Errorf("could not create location %s: %v", location.Name, err)
		return err
	}
	if !ok {
		result, err := query.Locations().Create(*location)
		if err != nil {
			s.ctx.Logger().Errorf("could not create location %s: %v", location.Name, err)
			return err
		}
		location = &result
	} else {
		result, ok, err := query.Locations().Update(queries.LocationsParams{
			ID: stringutil.ToPtr(currentLocation.ID),
		}, *location)
		if err != nil {
			s.ctx.Logger().Errorf("could not create location %s: %v", location.Name, err)
			return err
		}
		if !ok {
			s.ctx.Logger().Errorf("could not create location %s", location.Name)
			return nil
		}
		location = &result
	}
	s.ctx.Logger().Errorf("location %s successfully", location.Name)
	return nil
}
