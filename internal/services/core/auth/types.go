package auth

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"

	"github.com/golang-jwt/jwt/v4"
)

type (
	jwtClaims struct {
		*jwt.RegisteredClaims
		ID string `json:"id"`
	}
)

type (
	LoginData struct {
		Username  string `json:"username" example:"johndoe@example.com" validate:"required,email"`
		Password  string `json:"password" example:"******" validate:"required"`
		Remember  bool   `json:"remember" example:"true"`
		IPAddress string `json:"-" swaggerignore:"true"`
	}

	LoginResult struct {
		Code       types.ServiceResultCode
		Payload    *LoginResultPayload
		Validation types.ValidationResult
	}

	LoginResultPayload struct {
		User  entities.User `json:"user"`
		Token string        `json:"token"`
	}
)

type (
	JWTAuthData struct {
		Token string `json:"token"`
	}

	JWTAuthResult struct {
		Code    types.ServiceResultCode
		Payload *JWTAuthResultPayload
	}

	JWTAuthResultPayload struct {
		Token string        `json:"token"`
		User  entities.User `json:"user"`
	}
)
