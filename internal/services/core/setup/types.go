package setup

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
)

type (
	GetData struct {
		User *entities.User `json:"-"`
	}
	GetResult struct {
		Code    types.ServiceResultCode
		Payload GetResultPayload
	}

	GetResultPayload struct {
		Setup struct {
			Admin    bool               `json:"admin"`
			Database bool               `json:"database"`
			Location *entities.Location `json:"location"`
		} `json:"setup"`
	}
)

type (
	CreateData struct {
		Admin *struct {
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			EmailAddress string `json:"email_address"`
			PhoneNumber  string `json:"phone_number"`
			Password     string `json:"password"`
		} `json:"admin"`
		Location *struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"location"`
		Database *struct {
			SQLite *struct {
				EncryptionEnabled bool `json:"encryption_enabled"`
			} `json:"sqlite"`
			MariaDB *struct {
				Host         string `json:"host"`
				Port         uint32 `json:"port"`
				Username     string `json:"username"`
				Password     string `json:"password"`
				DatabaseName string `json:"database_name"`
			} `json:"mariadb"`
			Postgres *struct {
				Host         string `json:"host"`
				Port         uint32 `json:"port"`
				Username     string `json:"username"`
				Password     string `json:"password"`
				DatabaseName string `json:"database_name"`
			} `json:"postgres"`
		} `json:"database"`
	}
	CreateResult struct {
		Code       types.ServiceResultCode
		Payload    CreateResultPayload
		Validation types.ValidationResult
	}

	CreateResultPayload struct {
		Setup struct {
			Admin    bool `json:"admin"`
			Database bool `json:"database"`
			Location bool `json:"location"`
		} `json:"setup"`
	}
)
